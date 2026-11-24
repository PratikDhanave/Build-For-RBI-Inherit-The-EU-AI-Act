package consent

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidSubject       = errors.New("subject_id and category must be non-empty")
	ErrNonMonotonicTimeline = errors.New("consent event timestamp precedes latest event for subject/category")
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

func (l *Ledger) Grant(subjectID, category, reason string, at time.Time) (Event, error) {
	return l.append(subjectID, category, "GRANT", reason, at)
}

func (l *Ledger) Revoke(subjectID, category, reason string, at time.Time) (Event, error) {
	return l.append(subjectID, category, "REVOKE", reason, at)
}

func (l *Ledger) append(subjectID, category, action, reason string, at time.Time) (Event, error) {
	sub := strings.TrimSpace(subjectID)
	cat := strings.ToLower(strings.TrimSpace(category))
	if sub == "" || cat == "" {
		return Event{}, ErrInvalidSubject
	}
	ts := at.UTC()

	l.mu.Lock()
	defer l.mu.Unlock()
	for i := len(l.events) - 1; i >= 0; i-- {
		prev := l.events[i]
		if prev.SubjectID == sub && prev.Category == cat {
			if ts.Before(prev.Timestamp) {
				return Event{}, fmt.Errorf("%w: %s < %s", ErrNonMonotonicTimeline, ts.Format(time.RFC3339), prev.Timestamp.Format(time.RFC3339))
			}
			break
		}
	}
	ev := Event{
		Seq:       len(l.events) + 1,
		Timestamp: ts,
		SubjectID: sub,
		Category:  cat,
		Action:    action,
		Reason:    reason,
	}
	l.events = append(l.events, ev)
	return ev, nil
}

// IsLive returns whether the subject currently holds active consent for category
// along with the most recent ledger event for provable refusal reporting.
func (l *Ledger) IsLive(subjectID, category string) (bool, Event, bool) {
	sub := strings.TrimSpace(subjectID)
	cat := strings.ToLower(strings.TrimSpace(category))
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i := len(l.events) - 1; i >= 0; i-- {
		ev := l.events[i]
		if ev.SubjectID == sub && ev.Category == cat {
			return ev.Action == "GRANT", ev, true
		}
	}
	return false, Event{}, false
}

// History returns a copy of all immutable consent events for a subject.
func (l *Ledger) History(subjectID string) []Event {
	sub := strings.TrimSpace(subjectID)
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Event
	for _, ev := range l.events {
		if ev.SubjectID == sub {
			out = append(out, ev)
		}
	}
	return out
}
