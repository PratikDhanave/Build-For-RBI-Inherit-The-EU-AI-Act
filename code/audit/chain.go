package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/PratikDhanave/Build-For-RBI-Inherit-The-EU-AI-Act/code/gate"
)

// Entry represents one immutable row in the SHA-256 hash-chained audit log
// (RBI FREE-AI Rec 22 / EU AI Act Article 12).
type Entry struct {
	Seq       int    `json:"seq"`
	Timestamp string `json:"timestamp"`
	Actor     string `json:"actor"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Outcome   string `json:"outcome"`
	Reason    string `json:"reason,omitempty"`
	PrevHash  string `json:"prev_hash"`
	RowHash   string `json:"row_hash"`
}

// Manifest is emitted on export and always reports ChainVerified honestly —
// including when tampering has broken the chain.
type Manifest struct {
	Count         int    `json:"count"`
	ChainVerified bool   `json:"chain_verified"`
	BrokenAtSeq   int    `json:"broken_at_seq,omitempty"`
	BrokenField   string `json:"broken_field,omitempty"`
	PayloadSHA256 string `json:"sha256"`
	ExportedAt    string `json:"exported_at"`
}

// Log is an append-only, in-memory SHA-256 hash chain that also satisfies gate.AuditLogger.
type Log struct {
	mu      sync.Mutex
	Entries []Entry
	Now     func() time.Time
}

func NewLog() *Log {
	return &Log{
		Now: func() time.Time { return time.Now().UTC() },
	}
}

func computeRowHash(prev string, e Entry) string {
	type hashBody struct {
		Prev      string `json:"prev"`
		Seq       int    `json:"seq"`
		Timestamp string `json:"timestamp"`
		Actor     string `json:"actor"`
		Action    string `json:"action"`
		Target    string `json:"target"`
		Outcome   string `json:"outcome"`
		Reason    string `json:"reason,omitempty"`
		PrevHash  string `json:"prev_hash"`
	}
	b, _ := json.Marshal(hashBody{
		Prev:      prev,
		Seq:       e.Seq,
		Timestamp: e.Timestamp,
		Actor:     e.Actor,
		Action:    e.Action,
		Target:    e.Target,
		Outcome:   e.Outcome,
		Reason:    e.Reason,
		PrevHash:  e.PrevHash,
	})
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Append adds a new hash-chained entry bound to the previous row's digest.
func (l *Log) Append(actor, action, target, outcome, reason string) Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	prev := "genesis"
	if len(l.Entries) > 0 {
		prev = l.Entries[len(l.Entries)-1].RowHash
	}
	nowFn := l.Now
	if nowFn == nil {
		nowFn = func() time.Time { return time.Now().UTC() }
	}
	e := Entry{
		Seq:       len(l.Entries) + 1,
		Timestamp: nowFn().UTC().Format(time.RFC3339),
		Actor:     actor,
		Action:    action,
		Target:    target,
		Outcome:   outcome,
		Reason:    reason,
		PrevHash:  prev,
	}
	e.RowHash = computeRowHash(prev, e)
	l.Entries = append(l.Entries, e)
	return e
}

// RecordAllow implements gate.AuditLogger.
func (l *Log) RecordAllow(msg gate.Message) error {
	actor := msg.ActorID
	if actor == "" {
		actor = msg.ActorRole
	}
	l.Append(actor, "gate.dispatch", msg.AgentID, "ALLOW", "")
	return nil
}

// RecordDeny implements gate.AuditLogger.
func (l *Log) RecordDeny(msg gate.Message, policy, reason string) error {
	actor := msg.ActorID
	if actor == "" {
		actor = msg.ActorRole
	}
	l.Append(actor, "gate.dispatch", msg.AgentID, "DENY", fmt.Sprintf("[%s] %s", policy, reason))
	return nil
}

// Verify walks the entire chain from genesis and returns the first broken sequence/field if tampered.
func (l *Log) Verify() (ok bool, brokenSeq int, brokenField string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	prev := "genesis"
	for _, e := range l.Entries {
		if e.PrevHash != prev {
			return false, e.Seq, "prev_hash"
		}
		if computeRowHash(prev, e) != e.RowHash {
			return false, e.Seq, "row_hash"
		}
		prev = e.RowHash
	}
	return true, 0, ""
}

// ExportManifest produces a signed-ready summary manifest that never hides a broken chain.
func (l *Log) ExportManifest() Manifest {
	ok, seq, field := l.Verify()

	l.mu.Lock()
	defer l.mu.Unlock()

	lines := make([]string, 0, len(l.Entries))
	for _, e := range l.Entries {
		b, _ := json.Marshal(e)
		lines = append(lines, string(b))
	}
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	nowFn := l.Now
	if nowFn == nil {
		nowFn = func() time.Time { return time.Now().UTC() }
	}
	return Manifest{
		Count:         len(l.Entries),
		ChainVerified: ok,
		BrokenAtSeq:   seq,
		BrokenField:   field,
		PayloadSHA256: hex.EncodeToString(sum[:]),
		ExportedAt:    nowFn().UTC().Format(time.RFC3339),
	}
}
