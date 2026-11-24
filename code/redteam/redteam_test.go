package redteam

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/audit"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/config"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/gate"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/registry"
)

type corpusCase struct {
	ID           string `json:"id"`
	Category     string `json:"category"`
	OWASPASI     string `json:"owasp_asi"`
	AgentID      string `json:"agent_id"`
	Description  string `json:"description"`
	ActorRole    string `json:"actor_role"`
	DeclaredRisk string `json:"declared_risk"`
	ContainsPII  bool   `json:"contains_pii"`
	Region       string `json:"region"`
	IsAdvice     bool   `json:"is_advice"`
	Rationale    string `json:"rationale"`
	ConsentLive  bool   `json:"consent_live"`
	ExpectAllow  bool   `json:"expect_allow"`
	ExpectPolicy string `json:"expect_policy"`
}

// TestRedTeamCorpus executes the 20-case OWASP ASI-tagged adversarial corpus on every PR
// against the REAL board-approved policy.yaml + RegistryPolicy + ApprovalDriftPolicy
// (RBI FREE-AI Rec 20 / EU AI Act Article 15).
func TestRedTeamCorpus(t *testing.T) {
	cfg, err := config.LoadPolicyYAML(filepath.Join("..", "config", "policy.yaml"))
	if err != nil {
		t.Fatalf("LoadPolicyYAML: %v", err)
	}
	raw, err := os.ReadFile("corpus.json")
	if err != nil {
		t.Fatalf("read corpus.json: %v", err)
	}
	var cases []corpusCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("unmarshal corpus.json: %v", err)
	}
	if len(cases) < 20 {
		t.Fatalf("expected at least 20 red-team cases, got %d", len(cases))
	}

	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	validExp := now.Add(60 * 24 * time.Hour)

	reg := registry.NewLiveRegistry()
	reg.Register(registry.AgentDescriptor{ID: "educator", DeclaredRisk: "low"})
	reg.Register(registry.AgentDescriptor{ID: "payment_orchestrator", DeclaredRisk: "high"})
	reg.Register(registry.AgentDescriptor{ID: "var_calculator", DeclaredRisk: "high"})

	approvals := map[string]*gate.ApprovalRecord{
		"payment_orchestrator": {Status: "approved", RiskClass: "high", ExpiresAt: &validExp},
		"var_calculator":       {Status: "approved", RiskClass: "medium"}, // silent drift
	}

	auditLog := audit.NewLog()
	g := &gate.CompositeGate{
		Policies: []gate.Policy{
			gate.RegistryPolicy{IsRegistered: reg.Has},
			gate.RBACPolicy{HighRiskRoles: cfg.HighRiskRoles},
			gate.ApprovalDriftPolicy{Register: approvals, Now: func() time.Time { return now }},
			gate.ResidencyPolicy{AllowedPIIRegions: cfg.PIIAllowedRegions, DenyRegions: cfg.DenyRegions},
			gate.ExplainabilityPolicy{},
			gate.ConsentPolicy{},
		},
		Audit: auditLog,
	}

	for _, tc := range cases {
		t.Run(tc.ID+"_"+tc.Category, func(t *testing.T) {
			if tc.Category != "baseline" && !strings.HasPrefix(tc.OWASPASI, "ASI") {
				t.Fatalf("%s missing valid OWASP ASI tag: %q", tc.ID, tc.OWASPASI)
			}
			agentID := tc.AgentID
			if agentID == "" {
				agentID = "educator"
			}
			msg := gate.Message{
				ActorID:      "redteam-" + tc.ID,
				SubjectID:    "subject-" + tc.ID,
				ActorRole:    tc.ActorRole,
				AgentID:      agentID,
				Region:       tc.Region,
				ContainsPII:  tc.ContainsPII,
				IsAdvice:     tc.IsAdvice,
				Rationale:    tc.Rationale,
				ConsentLive:  tc.ConsentLive,
				DeclaredRisk: tc.DeclaredRisk,
			}
			err := g.Allow(context.Background(), msg)
			if tc.ExpectAllow && err != nil {
				t.Fatalf("%s (%s): expected ALLOW, got %v", tc.ID, tc.Description, err)
			}
			if !tc.ExpectAllow {
				if err == nil {
					t.Fatalf("%s (%s): expected DENY [%s], got ALLOW", tc.ID, tc.Description, tc.ExpectPolicy)
				}
				if tc.ExpectPolicy != "" && !strings.Contains(err.Error(), "["+tc.ExpectPolicy+"]") {
					t.Fatalf("%s (%s): expected policy [%s], got %v", tc.ID, tc.Description, tc.ExpectPolicy, err)
				}
			}
		})
	}

	if ok, seq, field := auditLog.Verify(); !ok {
		t.Fatalf("audit chain broken after red-team suite at seq=%d (%s)", seq, field)
	}
}
