package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/audit"
)

func TestGenerateDualEvidencePack(t *testing.T) {
	out := t.TempDir()
	policyPath := filepath.Join("..", "..", "config", "policy.yaml")
	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)

	trustedPub, err := GeneratePack(out, policyPath, now)
	if err != nil {
		t.Fatalf("GeneratePack failed: %v", err)
	}

	required := []string{
		"INDEX.md",
		"01-policy/policy.yaml",
		"02-inventory/inventory-export.json",
		"03-approvals/approval-drift-report.json",
		"04-disclosure/art50-verification.json",
		"05-audit/audit-chain.jsonl",
		"05-audit/signed-manifest.json",
		"06-incidents/annexure-vi-sample.json",
	}
	for _, rel := range required {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Fatalf("expected artefact %s in evidence pack: %v", rel, err)
		}
	}

	raw, err := os.ReadFile(filepath.Join(out, "05-audit", "signed-manifest.json"))
	if err != nil {
		t.Fatalf("read signed-manifest.json: %v", err)
	}
	var sm audit.SignedManifest
	if err := json.Unmarshal(raw, &sm); err != nil {
		t.Fatalf("unmarshal signed-manifest.json: %v", err)
	}
	if !sm.Manifest.ChainVerified || !audit.VerifySignedManifest(sm, trustedPub) {
		t.Fatalf("expected valid chain and pinned-key Ed25519 signature in evidence pack: %+v", sm)
	}
}
