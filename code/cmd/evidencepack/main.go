package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
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

type packAgent struct {
	id   string
	risk string
}

func (a packAgent) ID() string           { return a.id }
func (a packAgent) DeclaredRisk() string { return a.risk }

// GeneratePack builds a complete, self-contained dual evidence pack for both
// RBI FREE-AI and the EU AI Act directly from live offline execution, returning
// the Ed25519 public key used to sign the manifest (ephemeral in offline demo mode;
// pin a KMS/HSM public key in production).
func GeneratePack(outDir, policyPath string, now time.Time) (ed25519.PublicKey, error) {
	cfg, err := config.LoadPolicyYAML(policyPath)
	if err != nil {
		return nil, err
	}
	policyBytes, err := os.ReadFile(policyPath)
	if err != nil {
		return nil, err
	}

	dirs := []string{
		"01-policy",
		"02-inventory",
		"03-approvals",
		"04-disclosure",
		"05-audit",
		"06-incidents",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(outDir, d), 0o755); err != nil {
			return nil, err
		}
	}

	// 01-policy (RBI Rec 14 / EU AI Act Art 17)
	if err := os.WriteFile(filepath.Join(outDir, "01-policy", "policy.yaml"), policyBytes, 0o644); err != nil {
		return nil, err
	}

	// 02-inventory (RBI Rec 23 / EU AI Act Annex IV)
	reg := registry.NewLiveRegistry()
	reg.Register(registry.AgentDescriptor{ID: "payment_orchestrator", Domain: "payments", DeclaredRisk: "high", Owner: "payments@bank", VendorModel: "deterministic+llm", DataCategories: []string{"txn"}, DisclosureMode: "first-event"})
	reg.Register(registry.AgentDescriptor{ID: "var_calculator", Domain: "market-risk", DeclaredRisk: "high", Owner: "treasury@bank", VendorModel: "quant-llm", DataCategories: []string{"portfolio"}, DisclosureMode: "first-event"})
	reg.Register(registry.AgentDescriptor{ID: "educator", Domain: "retail", DeclaredRisk: "low", Owner: "cx@bank", VendorModel: "local-slm", DataCategories: []string{"faq"}, DisclosureMode: "first-event"})
	invJSON, err := reg.Export(now).ToJSON()
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "02-inventory", "inventory-export.json"), invJSON, 0o644); err != nil {
		return nil, err
	}

	// 03-approvals (RBI Rec 16, 17 / EU AI Act Art 9, 14)
	validExp := now.Add(60 * 24 * time.Hour)
	expired := now.Add(-24 * time.Hour)
	approvals := map[string]*gate.ApprovalRecord{
		"payment_orchestrator": {Status: "approved", RiskClass: "high", ExpiresAt: &validExp},
		"kyc_orchestrator":     {Status: "approved", RiskClass: "high", ExpiresAt: &expired},
		"aml_monitor":          {Status: "rejected", RiskClass: "high"},
		"var_calculator":       {Status: "approved", RiskClass: "medium"},
		"educator":             nil,
	}
	type driftRow struct {
		AgentID string `json:"agent_id"`
		Risk    string `json:"declared_risk"`
		Verdict string `json:"verdict"`
		Reason  string `json:"reason"`
	}
	var driftReport []driftRow
	for _, a := range []packAgent{
		{"payment_orchestrator", "high"},
		{"kyc_orchestrator", "high"},
		{"aml_monitor", "high"},
		{"var_calculator", "high"},
		{"educator", "low"},
	} {
		err := gate.CheckApproval(approvals, a, now)
		if err == nil {
			driftReport = append(driftReport, driftRow{a.ID(), a.DeclaredRisk(), "RUN", "live"})
		} else {
			driftReport = append(driftReport, driftRow{a.ID(), a.DeclaredRisk(), "STOP", err.Error()})
		}
	}
	driftJSON, err := json.MarshalIndent(driftReport, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "03-approvals", "approval-drift-report.json"), driftJSON, 0o644); err != nil {
		return nil, err
	}

	// 04-disclosure + 05-audit (RBI Rec 15, 18, 22, 25 / EU AI Act Art 12, 50)
	ledger := consent.NewLedger()
	if _, err := ledger.Grant("cust-active", "pii", "onboarding opt-in", now.Add(-5*24*time.Hour)); err != nil {
		return nil, err
	}
	if _, err := ledger.Grant("cust-revoked", "pii", "Monday grant", now.Add(-4*24*time.Hour)); err != nil {
		return nil, err
	}
	if _, err := ledger.Revoke("cust-revoked", "pii", "customer revoked on Tuesday", now.Add(-3*24*time.Hour)); err != nil {
		return nil, err
	}

	auditLog := audit.NewLog()
	auditLog.Now = func() time.Time { return now }
	policyGate := &gate.CompositeGate{
		Policies: []gate.Policy{
			gate.RegistryPolicy{IsRegistered: reg.Has},
			gate.RBACPolicy{HighRiskRoles: cfg.HighRiskRoles},
			gate.ApprovalDriftPolicy{Register: approvals, Now: func() time.Time { return now }},
			gate.ResidencyPolicy{AllowedPIIRegions: cfg.PIIAllowedRegions, DenyRegions: cfg.DenyRegions},
			gate.ExplainabilityPolicy{},
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
	var envelopes []disclosure.Envelope
	for _, msg := range []gate.Message{
		{ActorID: "u-1", SubjectID: "u-1", ActorRole: "user", AgentID: "educator", Region: "in", DeclaredRisk: "low"},
		{ActorID: "adv-7", SubjectID: "cust-active", ActorRole: "advisor", AgentID: "var_calculator", Region: "in", DeclaredRisk: "high"},
		{ActorID: "adv-7", SubjectID: "cust-active", ActorRole: "advisor", AgentID: "educator", Region: "us", ContainsPII: true, DeclaredRisk: "medium"},
		{ActorID: "adv-7", SubjectID: "cust-revoked", ActorRole: "advisor", AgentID: "educator", Region: "in", ContainsPII: true, DeclaredRisk: "medium"},
		{ActorID: "u-1", SubjectID: "u-1", ActorRole: "user", AgentID: "educator", Region: "in", DeclaredRisk: "low"},
	} {
		if err := policyGate.Allow(context.Background(), msg); err == nil {
			envelopes = append(envelopes, disc.Wrap(msg.ActorID, "dispatched"))
		}
	}
	discJSON, err := json.MarshalIndent(envelopes, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "04-disclosure", "art50-verification.json"), discJSON, 0o644); err != nil {
		return nil, err
	}

	jsonlBytes, err := auditLog.ExportJSONL()
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "05-audit", "audit-chain.jsonl"), jsonlBytes, 0o644); err != nil {
		return nil, err
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	signed, err := auditLog.ExportSignedManifest(priv, "rfc3161://tsa.bank.internal/"+now.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	signedJSON, err := json.MarshalIndent(signed, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "05-audit", "signed-manifest.json"), signedJSON, 0o644); err != nil {
		return nil, err
	}

	// 06-incidents (RBI Rec 22 Annexure VI / EU AI Act Art 73)
	rep := incident.NewAnnexureVI(
		"INC-2026-0042",
		"ClearRoute Reference Bank",
		"var_calculator",
		"high",
		"Governance Drift (ErrRiskChanged)",
		"HIGH",
		"var_calculator approved as medium risk; refactor elevated runtime declaration to high without renewed CRO sign-off.",
		"ApprovalDriftGuard (RBI Rec 17)",
		"None — blocked pre-dispatch by Composite Policy Gate",
		"Blocked pre-dispatch; re-approval workflow opened.",
		signed.Manifest.ChainVerified,
		signed.Manifest.PayloadSHA256,
		now,
	)
	incJSON, err := rep.ToJSON()
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "06-incidents", "annexure-vi-sample.json"), incJSON, 0o644); err != nil {
		return nil, err
	}

	indexMD := fmt.Sprintf(`# Dual Regulatory Evidence Pack (RBI FREE-AI + EU AI Act)

- **Generated At:** %s
- **Board Policy Version:** %s (approved by %s)
- **Audit Chain Verified:** %v (Ed25519 pinned-key signature valid: %v)
- **Ephemeral Demo Signing Public Key (replace with KMS/HSM in production):** `+"`%s`"+`

| Folder | Artefact | RBI FREE-AI | EU AI Act |
| --- | --- | --- | --- |
| 01-policy | policy.yaml | Rec 14 | Article 17 |
| 02-inventory | inventory-export.json | Rec 23 | Annex IV |
| 03-approvals | approval-drift-report.json | Rec 16, 17 | Articles 9, 14 |
| 04-disclosure | art50-verification.json | Rec 18, 25 | Article 50 |
| 05-audit | audit-chain.jsonl + signed-manifest.json | Rec 22 | Article 12 |
| 06-incidents | annexure-vi-sample.json | Rec 22 (Annexure VI) | Article 73 |
`, now.UTC().Format(time.RFC3339), cfg.Version, cfg.ApprovedBy, signed.Manifest.ChainVerified, audit.VerifySignedManifest(signed, pub), hex.EncodeToString(pub))

	if err := os.WriteFile(filepath.Join(outDir, "INDEX.md"), []byte(indexMD), 0o644); err != nil {
		return nil, err
	}
	return pub, nil
}

func main() {
	out := flag.String("out", filepath.Join(os.TempDir(), "rbi-eu-evidence-pack"), "Output directory for dual evidence pack")
	policy := flag.String("policy", filepath.Join("config", "policy.yaml"), "Path to board-approved policy.yaml")
	flag.Parse()

	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	if _, err := GeneratePack(*out, *policy, now); err != nil {
		fmt.Fprintf(os.Stderr, "evidencepack error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[OK] Dual RBI + EU AI Act evidence pack generated at: %s\n", *out)
}
