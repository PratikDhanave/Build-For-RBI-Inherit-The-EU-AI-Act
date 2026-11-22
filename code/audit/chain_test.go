package audit

import (
	"testing"
	"time"
)

func TestHashChainDetectsTamper(t *testing.T) {
	fixed := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
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

	manifest := l.ExportManifest()
	if manifest.ChainVerified {
		t.Fatalf("manifest must report ChainVerified=false when tampered")
	}
	if manifest.Count != 3 || manifest.BrokenAtSeq != 2 {
		t.Fatalf("unexpected manifest: %+v", manifest)
	}
}
