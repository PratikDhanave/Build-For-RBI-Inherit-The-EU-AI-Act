package gate

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type simpleAgent struct {
	id   string
	risk string
}

func (a simpleAgent) ID() string           { return a.id }
func (a simpleAgent) DeclaredRisk() string { return a.risk }

type memoryAudit struct {
	allows    int
	denies    []string
	failAllow bool
}

func (m *memoryAudit) RecordAllow(_ Message) error {
	if m.failAllow {
		return errors.New("audit disk full")
	}
	m.allows++
	return nil
}

func (m *memoryAudit) RecordDeny(_ Message, policy, _ string) error {
	m.denies = append(m.denies, policy)
	return nil
}

func TestCompositeGateAndSecurityHardening(t *testing.T) {
	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	validExp := now.Add(60 * 24 * time.Hour)
	registered := map[string]bool{
		"educator":             true,
		"payment_orchestrator": true,
		"var_calculator":       true,
	}
	approvals := map[string]*ApprovalRecord{
		"payment_orchestrator": {Status: "approved", RiskClass: "high", ExpiresAt: &validExp},
		"var_calculator":       {Status: "approved", RiskClass: "medium"}, // silent drift
	}
	consents := map[string]bool{
		"cust-live":    true,
		"cust-revoked": false,
	}

	auditSink := &memoryAudit{}
	g := &CompositeGate{
		Policies: []Policy{
			RegistryPolicy{IsRegistered: func(id string) bool { return registered[id] }},
			RBACPolicy{HighRiskRoles: map[string]bool{"advisor": true, "operator": true}},
			ApprovalDriftPolicy{Register: approvals, Now: func() time.Time { return now }},
			ResidencyPolicy{
				AllowedPIIRegions: map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true},
				DenyRegions:       map[string]bool{"us": true},
			},
			ExplainabilityPolicy{},
			ConsentPolicy{LedgerCheck: func(subjectID, _ string) (bool, string) {
				if consents[subjectID] {
					return true, ""
				}
				return false, "revoked on Tuesday"
			}},
		},
		Audit: auditSink,
	}

	cases := []struct {
		name       string
		msg        Message
		wantAllow  bool
		wantPolicy string
	}{
		{"ordinary question", Message{ActorID: "u-1", SubjectID: "u-1", ActorRole: "user", AgentID: "educator", DeclaredRisk: "low", Region: "in"}, true, ""},
		{"shadow unregistered agent blocked", Message{ActorID: "adv-7", SubjectID: "cust-live", ActorRole: "advisor", AgentID: "shadow_agent", DeclaredRisk: "low", Region: "in"}, false, "registry"},
		{"uppercase HIGH risk blocked for user", Message{ActorID: "u-1", ActorRole: "user", AgentID: "payment_orchestrator", DeclaredRisk: " HIGH ", Region: "in"}, false, "rbac"},
		{"unknown risk tier fails closed", Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "educator", DeclaredRisk: "critical", Region: "in"}, false, "rbac"},
		{"drifted var_calculator blocked at gate", Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "var_calculator", DeclaredRisk: "high", Region: "in"}, false, "approval_drift"},
		{"PAN to US model denied", Message{ActorID: "adv-7", SubjectID: "cust-live", ActorRole: "advisor", AgentID: "educator", DeclaredRisk: "medium", ContainsPII: true, Region: "US"}, false, "residency"},
		{"whitespace-only rationale blocked", Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "educator", DeclaredRisk: "medium", IsAdvice: true, Rationale: "   \t\n ", Region: "in"}, false, "explainability"},
		{"advisor blocked when customer SubjectID revoked", Message{ActorID: "adv-7", SubjectID: "cust-revoked", ActorRole: "advisor", AgentID: "educator", DeclaredRisk: "medium", ContainsPII: true, Region: "in"}, false, "consent"},
		{"authorized high-risk dispatch allowed", Message{ActorID: "op-1", SubjectID: "cust-live", ActorRole: "operator", AgentID: "payment_orchestrator", DeclaredRisk: "high", ContainsPII: true, Region: "IN"}, true, ""},
	}

	for _, tc := range cases {
		err := g.Allow(context.Background(), tc.msg)
		if tc.wantAllow && err != nil {
			t.Fatalf("%s: expected ALLOW, got %v", tc.name, err)
		}
		if !tc.wantAllow {
			if err == nil {
				t.Fatalf("%s: expected DENY [%s], got ALLOW", tc.name, tc.wantPolicy)
			}
			if !strings.Contains(err.Error(), "["+tc.wantPolicy+"]") {
				t.Fatalf("%s: expected policy [%s], got %v", tc.name, tc.wantPolicy, err)
			}
		}
	}

	if auditSink.allows != 2 || len(auditSink.denies) != 7 {
		t.Fatalf("unexpected audit totals: allows=%d denies=%d", auditSink.allows, len(auditSink.denies))
	}

	// Failing audit sink MUST fail closed on ALLOW
	auditSink.failAllow = true
	err := g.Allow(context.Background(), Message{ActorID: "u-1", ActorRole: "user", AgentID: "educator", DeclaredRisk: "low", Region: "in"})
	if err == nil || !strings.Contains(err.Error(), "[audit]") {
		t.Fatalf("expected fail-closed [audit] error when audit sink fails, got %v", err)
	}
}

func TestCheckApprovalDriftAndNormalization(t *testing.T) {
	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	validExp := now.Add(60 * 24 * time.Hour)
	expired := now.Add(-24 * time.Hour)

	register := map[string]*ApprovalRecord{
		"payment_orchestrator": {Status: "approved", RiskClass: "high", ExpiresAt: &validExp},
		"kyc_orchestrator":     {Status: "approved", RiskClass: "high", ExpiresAt: &expired},
		"aml_monitor":          {Status: "rejected", RiskClass: "high"},
		"var_calculator":       {Status: "approved", RiskClass: "medium"}, // silent drift!
		"educator":             nil,
	}

	if err := CheckApproval(register, simpleAgent{"payment_orchestrator", "HIGH"}, now); err != nil {
		t.Fatalf("expected payment_orchestrator allowed after case normalization, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"kyc_orchestrator", "high"}, now); !errors.Is(err, ErrExpired) {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"aml_monitor", "high"}, now); !errors.Is(err, ErrNotApproved) {
		t.Fatalf("expected ErrNotApproved, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"var_calculator", " HIGH "}, now); !errors.Is(err, ErrRiskChanged) {
		t.Fatalf("expected ErrRiskChanged for var_calculator even with padded uppercase HIGH, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"educator", "unknown-tier"}, now); !errors.Is(err, ErrUnknownRisk) {
		t.Fatalf("expected ErrUnknownRisk on unknown tier, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"educator", "low"}, now); err != nil {
		t.Fatalf("expected low-risk educator allowed, got %v", err)
	}
}
