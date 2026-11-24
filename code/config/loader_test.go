package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBoardPolicyYAML(t *testing.T) {
	cfg, err := LoadPolicyYAML("policy.yaml")
	if err != nil {
		t.Fatalf("LoadPolicyYAML failed: %v", err)
	}
	if cfg.ApprovedBy != "board-risk-committee@bank" {
		t.Fatalf("unexpected approved_by: %q", cfg.ApprovedBy)
	}
	for _, r := range []string{"in", "adgm", "difc", "sama"} {
		if !cfg.PIIAllowedRegions[r] {
			t.Fatalf("expected region %q in PIIAllowedRegions", r)
		}
	}
	if cfg.PIIAllowedRegions["us"] || !cfg.DenyRegions["us"] {
		t.Fatalf("us must be explicitly denied and not allowed for PII")
	}
	if !cfg.HighRiskRoles["advisor"] || !cfg.HighRiskRoles["operator"] {
		t.Fatalf("expected advisor and operator in HighRiskRoles: %+v", cfg.HighRiskRoles)
	}
	if !cfg.EmitDisclosureOnFirstEvent || cfg.AuditHashAlgorithm != "sha256" {
		t.Fatalf("unexpected disclosure/audit config: %+v", cfg)
	}
}

func TestYAMLSecurityInvariants(t *testing.T) {
	// 1. deny_regions overrides pii_allowed_regions and # inside quotes is preserved
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.yaml")
	content := `version: "2026.11#rc1"
approved_by: "board#risk@bank"
policies:
  residency:
    pii_allowed_regions:
      - "in"
      - "us"
    deny_regions:
      - "us"
  rbac:
    high_risk_allowed_roles:
      - "operator"
  audit:
    hash_algorithm: "sha256"
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadPolicyYAML(p)
	if err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
	if cfg.Version != "2026.11#rc1" || cfg.ApprovedBy != "board#risk@bank" {
		t.Fatalf("quoted # was wrongly stripped: version=%q approved_by=%q", cfg.Version, cfg.ApprovedBy)
	}
	if cfg.PIIAllowedRegions["us"] {
		t.Fatalf("explicit deny_regions must strip us from PIIAllowedRegions")
	}

	// 2. Reject weak hash algorithm (md5)
	weak := filepath.Join(dir, "weak.yaml")
	_ = os.WriteFile(weak, []byte("approved_by: \"b\"\npolicies:\n  residency:\n    pii_allowed_regions:\n      - \"in\"\n  rbac:\n    high_risk_allowed_roles:\n      - \"operator\"\n  audit:\n    hash_algorithm: \"md5\"\n"), 0o644)
	if _, err := LoadPolicyYAML(weak); err == nil {
		t.Fatalf("expected md5 hash_algorithm rejected")
	}

	// 3. Reject tab indentation
	tabbed := filepath.Join(dir, "tabbed.yaml")
	_ = os.WriteFile(tabbed, []byte("approved_by: \"b\"\n\tpolicies:\n"), 0o644)
	if _, err := LoadPolicyYAML(tabbed); err == nil {
		t.Fatalf("expected tab-indented YAML rejected")
	}
}
