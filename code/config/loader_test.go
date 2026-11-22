package config

import "testing"

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
	if cfg.PIIAllowedRegions["us"] {
		t.Fatalf("us must not be allowed for PII")
	}
	if !cfg.HighRiskRoles["advisor"] || !cfg.HighRiskRoles["operator"] {
		t.Fatalf("expected advisor and operator in HighRiskRoles: %+v", cfg.HighRiskRoles)
	}
	if !cfg.EmitDisclosureOnFirstEvent || cfg.AuditHashAlgorithm != "sha256" {
		t.Fatalf("unexpected disclosure/audit config: %+v", cfg)
	}
}
