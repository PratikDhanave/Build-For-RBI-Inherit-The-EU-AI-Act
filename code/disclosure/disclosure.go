package disclosure

import (
	"strings"
	"sync"
)

const (
	DefaultNotice      = "AI Disclosure (RBI FREE-AI Rec 18, 25 · EU AI Act Article 50): You are interacting with an automated AI system."
	DefaultMaxSessions = 100_000
)

// Envelope wraps an outgoing response with optional first-event AI disclosure.
type Envelope struct {
	SessionID string `json:"session_id"`
	Disclosed bool   `json:"ai_disclosed"`
	Notice    string `json:"disclosure_notice,omitempty"`
	Payload   string `json:"payload"`
}

// Tracker ensures every session emits the mandatory AI disclosure on its very
// first event (EU AI Act Article 50 — binding law since 2 August 2026) without
// spamming subsequent turns. Session tracking is bounded (FIFO eviction) to prevent
// unbounded memory growth under random session-ID floods, and empty session IDs
// always fail safe by emitting disclosure.
type Tracker struct {
	mu          sync.Mutex
	seen        map[string]bool
	order       []string
	MaxSessions int
	Notice      string
}

func NewTracker() *Tracker {
	return &Tracker{
		seen:        make(map[string]bool),
		MaxSessions: DefaultMaxSessions,
		Notice:      DefaultNotice,
	}
}

func (t *Tracker) Wrap(sessionID, payload string) Envelope {
	t.mu.Lock()
	defer t.mu.Unlock()

	notice := t.Notice
	if notice == "" {
		notice = DefaultNotice
	}
	sid := strings.TrimSpace(sessionID)
	// Fail safe: anonymous / empty session IDs always receive Art. 50 disclosure.
	if sid == "" {
		return Envelope{
			SessionID: sessionID,
			Disclosed: true,
			Notice:    notice,
			Payload:   payload,
		}
	}

	if !t.seen[sid] {
		limit := t.MaxSessions
		if limit <= 0 {
			limit = DefaultMaxSessions
		}
		if len(t.order) >= limit {
			oldest := t.order[0]
			t.order = t.order[1:]
			delete(t.seen, oldest)
		}
		t.seen[sid] = true
		t.order = append(t.order, sid)
		return Envelope{
			SessionID: sid,
			Disclosed: true,
			Notice:    notice,
			Payload:   payload,
		}
	}
	return Envelope{
		SessionID: sid,
		Disclosed: false,
		Payload:   payload,
	}
}
