package disclosure

import "sync"

const DefaultNotice = "AI Disclosure (RBI FREE-AI Rec 18, 25 · EU AI Act Article 50): You are interacting with an automated AI system."

// Envelope wraps an outgoing response with optional first-event AI disclosure.
type Envelope struct {
	SessionID string `json:"session_id"`
	Disclosed bool   `json:"ai_disclosed"`
	Notice    string `json:"disclosure_notice,omitempty"`
	Payload   string `json:"payload"`
}

// Tracker ensures every session emits the mandatory AI disclosure on its very
// first event (EU AI Act Article 50 — binding law since 2 August 2026) without
// spamming subsequent turns.
type Tracker struct {
	mu     sync.Mutex
	seen   map[string]bool
	Notice string
}

func NewTracker() *Tracker {
	return &Tracker{
		seen:   make(map[string]bool),
		Notice: DefaultNotice,
	}
}

func (t *Tracker) Wrap(sessionID, payload string) Envelope {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.seen[sessionID] {
		t.seen[sessionID] = true
		notice := t.Notice
		if notice == "" {
			notice = DefaultNotice
		}
		return Envelope{
			SessionID: sessionID,
			Disclosed: true,
			Notice:    notice,
			Payload:   payload,
		}
	}
	return Envelope{
		SessionID: sessionID,
		Disclosed: false,
		Payload:   payload,
	}
}
