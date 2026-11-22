package audit

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"
)

func TestHashChainDetectsTamperAndSignsHonestManifest(t *testing.T) {
	fixed := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	l := NewLog()
	l.Now = func() time.Time { return fixed }

	l.Append("system", "consent.grant", "user-1", "OK", "")
	l.Append("advisor-7", "portfolio.read", "user-1", "ALLOW", "")
	l.Append("payment_orchestrator", "refund.issue", "txn-9931", "ALLOW", "")

	if ok, seq, field := l.Verify(); !ok {
		t.Fatalf("expected fresh chain valid, got broken at seq=%d field=%s", seq, field)
	}

	// Simulate an operator quietly editing history on row 2
	l.Entries[1].Actor = "advisor-99"

	ok, seq, field := l.Verify()
	if ok || seq != 2 || field != "row_hash" {
		t.Fatalf("expected broken at seq=2 row_hash, got ok=%v seq=%d field=%s", ok, seq, field)
	}

	signed := l.ExportSignedManifest(priv, "rfc3161://tsa.bank.internal/2026-11-28T10:00:00Z")
	if !VerifySignedManifest(signed) {
		t.Fatalf("expected Ed25519 signature over manifest to verify")
	}
	if signed.Manifest.ChainVerified {
		t.Fatalf("signed manifest must honestly report ChainVerified=false when tampered")
	}
	if signed.Manifest.Count != 3 || signed.Manifest.BrokenAtSeq != 2 {
		t.Fatalf("unexpected manifest: %+v", signed.Manifest)
	}

	// Attempting to forge ChainVerified=true on the exported manifest breaks the Ed25519 signature
	signed.Manifest.ChainVerified = true
	if VerifySignedManifest(signed) {
		t.Fatalf("forged manifest must fail Ed25519 signature verification")
	}
}
