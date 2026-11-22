package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/audit"
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

	// 1. Live Registry Export (RBI Rec 23 / EU AI Act Annex IV)
	reg := registry.NewLiveRegistry()
	reg.Register(registry.AgentDescriptor{ID: "payment_orchestrator", Domain: "payments", DeclaredRisk: "high", Owner: "payments@bank", VendorModel: "deterministic+llm", DataCategories: []string{"txn"}, DisclosureMode: "first-event"})
	reg.Register(registry.AgentDescriptor{ID: "var_calculator", Domain: "market-risk", DeclaredRisk: "high", Owner: "treasury@bank", VendorModel: "quant-llm", DataCategories: []string{"portfolio"}, DisclosureMode: "first-event"})
	reg.Register(registry.AgentDescriptor{ID: "educator", Domain: "retail", DeclaredRisk: "low", Owner: "cx@bank", VendorModel: "local-slm", DataCategories: []string{"faq"}, DisclosureMode: "first-event"})
	inv := reg.Export(now)
	fmt.Printf("\n[1] Live Registry Inventory (Rec 23 / Annex IV): %d live agents (zero-drift)\n", inv.TotalLiveAgents)

	// 2. Composite Policy Gate + Hash-Chained Audit Log (Rec 15, 16, 19, 22 / Art 12)
	auditLog := audit.NewLog()
	auditLog.Now = func() time.Time { return now }
	policyGate := &gate.CompositeGate{
		Policies: []gate.Policy{
			gate.ResidencyPolicy{AllowedPIIRegions: map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true}},
			gate.ExplainabilityPolicy{},
			gate.RBACPolicy{HighRiskRoles: map[string]bool{"advisor": true, "operator": true}},
			gate.ConsentPolicy{},
		},
		Audit: auditLog,
	}

	fmt.Println("\n[2] Composite Policy Gate (Pre-Dispatch):")
	samples := []struct {
		label string
		msg   gate.Message
	}{
		{"ordinary question", gate.Message{ActorID: "u-1", ActorRole: "user", AgentID: "educator", Region: "in", DeclaredRisk: "low", ConsentLive: true}},
		{"PAN -> US model", gate.Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "var_calculator", Region: "us", ContainsPII: true, DeclaredRisk: "medium", ConsentLive: true}},
		{"advice with no reason", gate.Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "var_calculator", Region: "in", IsAdvice: true, Rationale: "", DeclaredRisk: "medium", ConsentLive: true}},
		{"PII, consent revoked", gate.Message{ActorID: "adv-7", ActorRole: "advisor", AgentID: "var_calculator", Region: "in", ContainsPII: true, DeclaredRisk: "medium", ConsentLive: false}},
		{"balance lookup", gate.Message{ActorID: "u-1", ActorRole: "user", AgentID: "educator", Region: "in", DeclaredRisk: "low", ConsentLive: true}},
	}
	for _, s := range samples {
		err := policyGate.Allow(context.Background(), s.msg)
		if err == nil {
			fmt.Printf("    ALLOW  %-24s -> dispatched\n", s.label)
		} else {
			fmt.Printf("    %s  (%s)\n", err.Error(), s.label)
		}
	}

	// 3. Finding 1 — Runtime Approval Drift Guard (RBI Rec 17)
	fmt.Println("\n[3] Finding 1 — Runtime Risk-Drift Guard (Rec 17):")
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

	// 4. Finding 2 — Tamper-Evident Audit Export & Annexure VI Incident Report (Rec 22 / Art 12 & 73)
	fmt.Println("\n[4] Finding 2 — SHA-256 Audit Chain & Tamper Detection (Rec 22 / Art 12):")
	cleanManifest := auditLog.ExportManifest()
	fmt.Printf("    fresh chain            chain_verified=%v (entries=%d, sha256=%s…)\n",
		cleanManifest.ChainVerified, cleanManifest.Count, cleanManifest.PayloadSHA256[:16])

	// Quietly rewrite entry 2
	auditLog.Entries[1].Actor = "advisor-99"
	tamperedManifest := auditLog.ExportManifest()
	fmt.Printf("    after editing row 2    chain_verified=%v (broken_at_seq=%d, field=%s)\n",
		tamperedManifest.ChainVerified, tamperedManifest.BrokenAtSeq, tamperedManifest.BrokenField)

	fmt.Println("\n[5] RBI Annexure VI / EU AI Act Article 73 Incident Artefact:")
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
		cleanManifest.ChainVerified,
		cleanManifest.PayloadSHA256[:24],
		now,
	)
	b, _ := json.MarshalIndent(rep, "    ", "  ")
	fmt.Println("   ", string(b))
}
