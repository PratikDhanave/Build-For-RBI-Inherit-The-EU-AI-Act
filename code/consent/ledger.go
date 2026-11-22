package consent

import (
	"sync"
	"time"
)

// Event is an immutable grant or revocation entry in the append-only consent ledger.
// Past entries are never overwritten or deleted so auditors can reconstruct
// why an in-flight request failed after a Tuesday revocation (RBI Rec 15).
type Event struct {
	Seq       int       `json:"seq"`
	Timestamp time.Time `json:"timestamp"`
	SubjectID string    `json:"subject_id"`
	Category  string    `json:"category"`
	Action    string    `json:"action"` // "GRANT" or "REVOKE"
	Reason    string    `json:"reason,omitempty"`
}

// Ledger stores consent state strictly as an append-only event history.
type Ledger struct {
	mu     sync.RWMutex
	events []Event
}

func NewLedger() *Ledger {
	return &Ledger{}
}

func (l *Ledger) Grant(subjectID, category, reason string, at time.Time) Event {
	return l.append(subjectID, category, "GRANT", reason, at)
}

func (l *Ledger) Revoke(subjectID, category, reason string, at time.Time) Event {
	return l.append(subjectID, category, "REVOKE", reason, at)
}

func (l *Ledger) append(subjectID, category, action, reason string, at time.Time) Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	ev := Event{
		Seq:       len(l.events) + 1,
		Timestamp: at.UTC(),
		SubjectID: subjectID,
		Category:  category,
		Action:    action,
		Reason:    reason,
	}
	l.events = append(l.events, ev)
	return ev
}

// IsLive returns whether the subject currently holds active consent for category
// along with the most recent ledger event for provable refusal reporting.
func (l *Ledger) IsLive(subjectID, category string) (bool, Event, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i := len(l.events) - 1; i >= 0; i-- {
		ev := l.events[i]
		if ev.SubjectID == subjectID && ev.Category == category {
			return ev.Action == "GRANT", ev, true
		}
	}
	return false, Event{}, false
}

// History returns a copy of all immutable consent events for a subject.
func (l *Ledger) History(subjectID string) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Event
	for _, ev := range l.events {
		if ev.SubjectID == subjectID {
			out = append(out, ev)
		}
	}
	return out
}
