package disclosure

import "testing"

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
}
