package consent

import (
	"errors"
	"testing"
	"time"
)

func TestAppendOnlyConsentRevocationOnTuesday(t *testing.T) {
	mon := time.Date(2026, 11, 23, 9, 0, 0, 0, time.UTC)
	tue := time.Date(2026, 11, 24, 14, 30, 0, 0, time.UTC)

	ledger := NewLedger()
	if _, err := ledger.Grant("user-1", "pii", "onboarding opt-in", mon); err != nil {
		t.Fatalf("Grant failed: %v", err)
	}

	live, ev, ok := ledger.IsLive("user-1", "pii")
	if !ok || !live || ev.Action != "GRANT" {
		t.Fatalf("expected live consent on Monday, got live=%v ev=%+v", live, ev)
	}

	// Customer revokes on Tuesday: original GRANT remains intact in ledger
	if _, err := ledger.Revoke("user-1", "pii", "customer revoked on Tuesday", tue); err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}

	live, ev, ok = ledger.IsLive("user-1", "pii")
	if !ok || live || ev.Action != "REVOKE" {
		t.Fatalf("expected revoked consent on Tuesday, got live=%v ev=%+v", live, ev)
	}

	// Backdated grant cannot silently override Tuesday revocation
	if _, err := ledger.Grant("user-1", "pii", "backdated replay", mon); !errors.Is(err, ErrNonMonotonicTimeline) {
		t.Fatalf("expected ErrNonMonotonicTimeline on backdated event, got %v", err)
	}

	// Empty subject rejected
	if _, err := ledger.Grant("", "pii", "empty", tue); !errors.Is(err, ErrInvalidSubject) {
		t.Fatalf("expected ErrInvalidSubject on empty subject, got %v", err)
	}

	hist := ledger.History("user-1")
	if len(hist) != 2 || hist[0].Action != "GRANT" || hist[1].Action != "REVOKE" {
		t.Fatalf("append-only history must preserve both events, got %+v", hist)
	}
}
