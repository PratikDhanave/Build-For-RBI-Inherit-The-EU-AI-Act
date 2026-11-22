package gate

import (
	"context"
	"errors"
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
	allows int
	denies []string
}

func (m *memoryAudit) RecordAllow(_ Message) error {
	m.allows++
	return nil
}

func (m *memoryAudit) RecordDeny(_ Message, policy, _ string) error {
	m.denies = append(m.denies, policy)
	return nil
}

func TestCompositeGate(t *testing.T) {
	audit := &memoryAudit{}
	g := &CompositeGate{
		Policies: []Policy{
			ResidencyPolicy{AllowedPIIRegions: map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true}},
			ExplainabilityPolicy{},
			RBACPolicy{HighRiskRoles: map[string]bool{"advisor": true, "operator": true}},
			ConsentPolicy{},
		},
		Audit: audit,
	}

	cases := []struct {
		name      string
		msg       Message
		wantAllow bool
	}{
		{"ordinary question", Message{ActorRole: "user", DeclaredRisk: "low", Region: "in", ConsentLive: true}, true},
		{"PAN to US model", Message{ActorRole: "advisor", DeclaredRisk: "medium", ContainsPII: true, Region: "us", ConsentLive: true}, false},
		{"advice without rationale", Message{ActorRole: "advisor", DeclaredRisk: "medium", IsAdvice: true, Rationale: "", ConsentLive: true}, false},
		{"user hits high agent", Message{ActorRole: "user", DeclaredRisk: "high", Region: "in", ConsentLive: true}, false},
		{"PII with revoked consent", Message{ActorRole: "advisor", DeclaredRisk: "medium", ContainsPII: true, Region: "in", ConsentLive: false}, false},
		{"balance lookup", Message{ActorRole: "user", DeclaredRisk: "low", Region: "in", ConsentLive: true}, true},
	}

	for _, tc := range cases {
		err := g.Allow(context.Background(), tc.msg)
		if tc.wantAllow && err != nil {
			t.Fatalf("%s: expected ALLOW, got %v", tc.name, err)
		}
		if !tc.wantAllow && err == nil {
			t.Fatalf("%s: expected DENY, got ALLOW", tc.name)
		}
	}

	if audit.allows != 2 || len(audit.denies) != 4 {
		t.Fatalf("unexpected audit totals: allows=%d denies=%d", audit.allows, len(audit.denies))
	}
}

func TestCheckApprovalDrift(t *testing.T) {
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

	if err := CheckApproval(register, simpleAgent{"payment_orchestrator", "high"}, now); err != nil {
		t.Fatalf("expected payment_orchestrator allowed, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"kyc_orchestrator", "high"}, now); !errors.Is(err, ErrExpired) {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"aml_monitor", "high"}, now); !errors.Is(err, ErrNotApproved) {
		t.Fatalf("expected ErrNotApproved, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"var_calculator", "high"}, now); !errors.Is(err, ErrRiskChanged) {
		t.Fatalf("expected ErrRiskChanged for var_calculator, got %v", err)
	}
	if err := CheckApproval(register, simpleAgent{"educator", "low"}, now); err != nil {
		t.Fatalf("expected low-risk educator allowed, got %v", err)
	}
}
