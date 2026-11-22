package redteam

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/audit"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/gate"
)

type corpusCase struct {
	ID           string `json:"id"`
	Category     string `json:"category"`
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

// TestRedTeamCorpus executes the adversarial corpus on every PR
// (RBI FREE-AI Rec 20 / EU AI Act Article 15).
func TestRedTeamCorpus(t *testing.T) {
	raw, err := os.ReadFile("corpus.json")
	if err != nil {
		t.Fatalf("read corpus.json: %v", err)
	}
	var cases []corpusCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("unmarshal corpus.json: %v", err)
	}
	if len(cases) < 16 {
		t.Fatalf("expected at least 16 red-team cases, got %d", len(cases))
	}

	auditLog := audit.NewLog()
	g := &gate.CompositeGate{
		Policies: []gate.Policy{
			gate.ResidencyPolicy{AllowedPIIRegions: map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true}},
			gate.ExplainabilityPolicy{},
			gate.RBACPolicy{HighRiskRoles: map[string]bool{"advisor": true, "operator": true}},
			gate.ConsentPolicy{},
		},
		Audit: auditLog,
	}

	for _, tc := range cases {
		t.Run(tc.ID+"_"+tc.Category, func(t *testing.T) {
			msg := gate.Message{
				ActorID:      "redteam-" + tc.ID,
				ActorRole:    tc.ActorRole,
				AgentID:      "target_agent",
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
