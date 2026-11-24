package audit

import (
	"crypto/ed25519"
	"crypto/rand"
	"sync"
	"testing"
	"time"
)

func TestHashChainDetectsTamperAndEnforcesPinnedKeyCrypto(t *testing.T) {
	fixed := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	trustedPub, trustedPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	_, attackerPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("attacker GenerateKey failed: %v", err)
	}

	l := NewLog()
	l.Now = func() time.Time { return fixed }

	// Empty / truncated log must NOT report ChainVerified=true
	if ok, _, field := l.Verify(); ok || field != "empty_log" {
		t.Fatalf("expected empty log to fail verify with empty_log, got ok=%v field=%s", ok, field)
	}

	l.Append("system", "consent.grant", "user-1", "OK", "")
	l.Append("advisor-7", "portfolio.read", "user-1", "ALLOW", "")
	l.Append("payment_orchestrator", "refund.issue", "txn-9931", "ALLOW", "")

	if ok, seq, field := l.Verify(); !ok {
		t.Fatalf("expected fresh chain valid, got broken at seq=%d field=%s", seq, field)
	}

	// Simulate an operator quietly editing history on row 2
	if err := l.TamperRow(1, func(e *Entry) { e.Actor = "advisor-99" }); err != nil {
		t.Fatalf("TamperRow failed: %v", err)
	}

	ok, seq, field := l.Verify()
	if ok || seq != 2 || field != "row_hash" {
		t.Fatalf("expected broken at seq=2 row_hash, got ok=%v seq=%d field=%s", ok, seq, field)
	}

	signed, err := l.ExportSignedManifest(trustedPriv, "rfc3161://tsa.bank.internal/2026-11-28T10:00:00Z")
	if err != nil {
		t.Fatalf("ExportSignedManifest failed: %v", err)
	}
	if !VerifySignedManifest(signed, trustedPub) {
		t.Fatalf("expected Ed25519 signature to verify against trusted public key")
	}
	if signed.Manifest.ChainVerified {
		t.Fatalf("signed manifest must honestly report ChainVerified=false when tampered")
	}
	if signed.Manifest.Count != 3 || signed.Manifest.BrokenAtSeq != 2 {
		t.Fatalf("unexpected manifest: %+v", signed.Manifest)
	}

	// Attack 1: Forging ChainVerified=true breaks signature
	forgedChain := signed
	forgedChain.Manifest.ChainVerified = true
	if VerifySignedManifest(forgedChain, trustedPub) {
		t.Fatalf("forged ChainVerified=true must fail verification")
	}

	// Attack 2: Backdating TimestampAnchor breaks signature
	forgedAnchor := signed
	forgedAnchor.TimestampAnchor = "rfc3161://attacker.evil/2020-01-01T00:00:00Z"
	if VerifySignedManifest(forgedAnchor, trustedPub) {
		t.Fatalf("tampered TimestampAnchor must fail verification")
	}

	// Attack 3: Self-signed key-substitution attack must fail against caller-pinned trustedPub
	attackerSigned, err := l.ExportSignedManifest(attackerPriv, "rfc3161://tsa.bank.internal/2026-11-28T10:00:00Z")
	if err != nil {
		t.Fatalf("attacker export failed: %v", err)
	}
	if VerifySignedManifest(attackerSigned, trustedPub) {
		t.Fatalf("self-signed key-substitution manifest MUST fail against pinned trustedPub")
	}
}

func TestConcurrentAppendAndExportRaceSafe(t *testing.T) {
	l := NewLog()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				l.Append("worker", "gate.dispatch", "educator", "ALLOW", "")
				_, _ = l.ExportManifest()
			}
		}()
	}
	wg.Wait()

	if ok, seq, field := l.Verify(); !ok {
		t.Fatalf("concurrent chain broken at seq=%d (%s)", seq, field)
	}
}
