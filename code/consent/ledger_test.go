package consent

import (
	"testing"
	"time"
)

func TestAppendOnlyConsentRevocationOnTuesday(t *testing.T) {
	mon := time.Date(2026, 11, 23, 9, 0, 0, 0, time.UTC)
	tue := time.Date(2026, 11, 24, 14, 30, 0, 0, time.UTC)

	ledger := NewLedger()
	ledger.Grant("user-1", "pii", "onboarding opt-in", mon)

	live, ev, ok := ledger.IsLive("user-1", "pii")
	if !ok || !live || ev.Action != "GRANT" {
		t.Fatalf("expected live consent on Monday, got live=%v ev=%+v", live, ev)
	}

	// Customer revokes on Tuesday: original GRANT remains intact in ledger
	ledger.Revoke("user-1", "pii", "customer revoked on Tuesday", tue)

	live, ev, ok = ledger.IsLive("user-1", "pii")
	if !ok || live || ev.Action != "REVOKE" {
		t.Fatalf("expected revoked consent on Tuesday, got live=%v ev=%+v", live, ev)
	}

	hist := ledger.History("user-1")
	if len(hist) != 2 || hist[0].Action != "GRANT" || hist[1].Action != "REVOKE" {
		t.Fatalf("append-only history must preserve both events, got %+v", hist)
	}
}
