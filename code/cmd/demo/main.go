package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/audit"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/config"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/consent"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/disclosure"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/gate"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/incident"
	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/registry"
)

type demoAgent struct {
	id   string
	risk string
}

func (a demoAgent) ID() string           { return a.id }
func (a demoAgent) DeclaredRisk() string { return a.risk }

func main() {
	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)

	fmt.Println("====================================================================")
	fmt.Println(" Build for RBI, Inherit the EU AI Act — End-to-End Offline Demo")
	fmt.Println(" Zero cloud account · Zero API keys · Zero network")
	fmt.Println("====================================================================")

	// 1. Load Board-Approved Policy YAML (RBI Rec 14 / EU AI Act Art 17)
	policyPath := filepath.Join("config", "policy.yaml")
	cfg, err := config.LoadPolicyYAML(policyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed loading %s: %v\n", policyPath, err)
		os.Exit(1)
	}
	fmt.Printf("\n[1] Board Policy YAML Loaded (Rec 14 / Art 17): v%s approved by %s\n", cfg.Version, cfg.ApprovedBy)

	// 2. Live Registry Export (RBI Rec 23 / EU AI Act Annex IV)
	reg := registry.NewLiveRegistry()
	reg.Register(registry.AgentDescriptor{ID: "payment_orchestrator", Domain: "payments", DeclaredRisk: "high", Owner: "payments@bank", VendorModel: "deterministic+llm", DataCategories: []string{"txn"}, DisclosureMode: "first-event"})
	reg.Register(registry.AgentDescriptor{ID: "var_calculator", Domain: "market-risk", DeclaredRisk: "high", Owner: "treasury@bank", VendorModel: "quant-llm", DataCategories: []string{"portfolio"}, DisclosureMode: "first-event"})
	reg.Register(registry.AgentDescriptor{ID: "educator", Domain: "retail", DeclaredRisk: "low", Owner: "cx@bank", VendorModel: "local-slm", DataCategories: []string{"faq"}, DisclosureMode: "first-event"})
	inv := reg.Export(now)
	fmt.Printf("[2] Live Registry Inventory (Rec 23 / Annex IV): %d live agents (zero-drift)\n", inv.TotalLiveAgents)

	// 3. Append-Only Consent Ledger (RBI Rec 15 — Tuesday Revocation)
	ledger := consent.NewLedger()
	ledger.Grant("adv-7", "pii", "onboarding opt-in", now.Add(-5*24*time.Hour))
	ledger.Grant("user-revoked", "pii", "Monday grant", now.Add(-4*24*time.Hour))
	ledger.Revoke("user-revoked", "pii", "customer revoked on Tuesday", now.Add(-3*24*time.Hour))

	// 4. Composite Policy Gate + First-Event AI Disclosure (Rec 15, 16, 18, 19, 25 / Art 50)
	auditLog := audit.NewLog()
	auditLog.Now = func() time.Time { return now }
	policyGate := &gate.CompositeGate{
		Policies: []gate.Policy{
			gate.ResidencyPolicy{AllowedPIIRegions: cfg.PIIAllowedRegions},
			gate.ExplainabilityPolicy{},
			gate.RBACPolicy{HighRiskRoles: cfg.HighRiskRoles},
			gate.ConsentPolicy{LedgerCheck: func(subjectID, category string) (bool, string) {
				live, ev, ok := ledger.IsLive(subjectID, category)
				if !ok {
					return false, "no consent record"
				}
				return live, ev.Reason
			}},
		},
		Audit: auditLog,
	}
	disc := disclosure.NewTracker()

	fmt.Println("\n[3] Composite Policy Gate + Append-Only Consent + First-Event Disclosure:")
	samples := []struct {
		label string
		msg   gate.Message
	}{
		{"ordinary question", gate.Message{ActorID: "u-1", ActorRole: "user", AgentID: "educator", Region: "in", DeclaredRisk: "low"}},
		{"PAN -> US model", gate.Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "var_calculator", Region: "us", ContainsPII: true, DeclaredRisk: "medium"}},
		{"advice with no reason", gate.Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "var_calculator", Region: "in", IsAdvice: true, Rationale: "", DeclaredRisk: "medium"}},
		{"PII, consent revoked", gate.Message{ActorID: "user-revoked", ActorRole: "advisor", AgentID: "var_calculator", Region: "in", ContainsPII: true, DeclaredRisk: "medium"}},
		{"balance lookup", gate.Message{ActorID: "u-1", ActorRole: "user", AgentID: "educator", Region: "in", DeclaredRisk: "low"}},
	}
	for _, s := range samples {
		err := policyGate.Allow(context.Background(), s.msg)
		if err == nil {
			env := disc.Wrap(s.msg.ActorID, "ok")
			fmt.Printf("    ALLOW  %-24s -> dispatched (art50_disclosed=%v)\n", s.label, env.Disclosed)
		} else {
			fmt.Printf("    %s  (%s)\n", err.Error(), s.label)
		}
	}

	// 5. Finding 1 — Runtime Approval Drift Guard (RBI Rec 17)
	fmt.Println("\n[4] Finding 1 — Runtime Risk-Drift Guard (Rec 17):")
	validExp := now.Add(60 * 24 * time.Hour)
	expired := now.Add(-24 * time.Hour)
	approvals := map[string]*gate.ApprovalRecord{
		"payment_orchestrator": {Status: "approved", RiskClass: "high", ExpiresAt: &validExp},
		"kyc_orchestrator":     {Status: "approved", RiskClass: "high", ExpiresAt: &expired},
		"aml_monitor":          {Status: "rejected", RiskClass: "high"},
		"var_calculator":       {Status: "approved", RiskClass: "medium"},
		"educator":             nil,
	}
	agents := []demoAgent{
		{"payment_orchestrator", "high"},
		{"kyc_orchestrator", "high"},
		{"aml_monitor", "high"},
		{"var_calculator", "high"},
		{"educator", "low"},
	}
	for _, a := range agents {
		if err := gate.CheckApproval(approvals, a, now); err != nil {
			fmt.Printf("    STOP   %-22s %v\n", a.ID(), err)
		} else {
			fmt.Printf("    RUN    %-22s live\n", a.ID())
		}
	}

	// 6. Finding 2 — SHA-256 Audit Chain + Ed25519 Signed Manifest (Rec 22 / Art 12)
	fmt.Println("\n[5] Finding 2 — SHA-256 Chain + Ed25519 Signed Manifest (Rec 22 / Art 12):")
	_, priv, _ := ed25519.GenerateKey(rand.Reader)
	cleanSigned := auditLog.ExportSignedManifest(priv, "rfc3161://tsa.bank.internal/2026-11-28T10:00:00Z")
	fmt.Printf("    fresh chain            chain_verified=%v sig_valid=%v (entries=%d)\n",
		cleanSigned.Manifest.ChainVerified, audit.VerifySignedManifest(cleanSigned), cleanSigned.Manifest.Count)

	// Quietly rewrite entry 2
	auditLog.Entries[1].Actor = "advisor-99"
	tamperedSigned := auditLog.ExportSignedManifest(priv, "rfc3161://tsa.bank.internal/2026-11-28T10:00:01Z")
	fmt.Printf("    after editing row 2    chain_verified=%v sig_valid=%v (broken_at_seq=%d, field=%s)\n",
		tamperedSigned.Manifest.ChainVerified, audit.VerifySignedManifest(tamperedSigned),
		tamperedSigned.Manifest.BrokenAtSeq, tamperedSigned.Manifest.BrokenField)

	// 7. RBI Annexure VI / EU AI Act Article 73 Incident Report
	fmt.Println("\n[6] RBI Annexure VI / EU AI Act Article 73 Incident Artefact:")
	rep := incident.NewAnnexureVI(
		"INC-2026-0042",
		"ClearRoute Reference Bank",
		"var_calculator",
		"high",
		"Governance Drift (ErrRiskChanged)",
		"HIGH",
		"var_calculator approved at medium risk; refactor elevated runtime declaration to high without renewed approval.",
		"ApprovalDriftGuard (RBI Rec 17)",
		"Blocked pre-dispatch; CRO re-approval workflow triggered.",
		cleanSigned.Manifest.ChainVerified,
		cleanSigned.Manifest.PayloadSHA256[:24],
		now,
	)
	b, _ := json.MarshalIndent(rep, "    ", "  ")
	fmt.Println("   ", string(b))
}
