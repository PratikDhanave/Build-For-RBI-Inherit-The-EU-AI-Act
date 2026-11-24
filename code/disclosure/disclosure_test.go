package disclosure

import (
	"fmt"
	"testing"
)

func TestFirstEventDisclosureOnlyOncePerSession(t *testing.T) {
	tr := NewTracker()

	first := tr.Wrap("sess-101", "Your current balance is ₹1,42,000.")
	if !first.Disclosed || first.Notice == "" {
		t.Fatalf("expected AI disclosure on first event, got %+v", first)
	}

	second := tr.Wrap("sess-101", "Would you like a mini statement?")
	if second.Disclosed || second.Notice != "" {
		t.Fatalf("expected no duplicate disclosure on second event, got %+v", second)
	}

	other := tr.Wrap("sess-202", "Hello from ClearRoute assistant.")
	if !other.Disclosed {
		t.Fatalf("expected disclosure on first event of new session")
	}

	// Empty/anonymous session IDs must ALWAYS fail safe with disclosure
	anon1 := tr.Wrap("", "anonymous turn 1")
	anon2 := tr.Wrap("   ", "anonymous turn 2")
	if !anon1.Disclosed || !anon2.Disclosed {
		t.Fatalf("empty session IDs must always emit Art. 50 disclosure: %+v %+v", anon1, anon2)
	}
}

func TestBoundedSessionEvictionPreventsOOM(t *testing.T) {
	tr := NewTracker()
	tr.MaxSessions = 3
	for i := 0; i < 10; i++ {
		tr.Wrap(fmt.Sprintf("s-%d", i), "hi")
	}
	if len(tr.seen) > 3 || len(tr.order) > 3 {
		t.Fatalf("session tracker exceeded MaxSessions bound: seen=%d order=%d", len(tr.seen), len(tr.order))
	}
}
