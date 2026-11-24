package audit

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
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
// including when tampering or truncation has broken the chain.
type Manifest struct {
	Count         int    `json:"count"`
	ChainVerified bool   `json:"chain_verified"`
	BrokenAtSeq   int    `json:"broken_at_seq,omitempty"`
	BrokenField   string `json:"broken_field,omitempty"`
	PayloadSHA256 string `json:"sha256"`
	ExportedAt    string `json:"exported_at"`
}

// Log is an append-only, concurrency-safe SHA-256 hash chain implementing gate.AuditLogger.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	Now     func() time.Time
}

func NewLog() *Log {
	return &Log{
		Now: func() time.Time { return time.Now().UTC() },
	}
}

func computeRowHash(prev string, e Entry) (string, error) {
	// Fields are declared in alphabetical JSON key order for canonical determinism.
	type canonicalRow struct {
		Action    string `json:"action"`
		Actor     string `json:"actor"`
		Outcome   string `json:"outcome"`
		PrevHash  string `json:"prev_hash"`
		Reason    string `json:"reason,omitempty"`
		Seq       int    `json:"seq"`
		Target    string `json:"target"`
		Timestamp string `json:"timestamp"`
	}
	b, err := json.Marshal(canonicalRow{
		Action:    e.Action,
		Actor:     e.Actor,
		Outcome:   e.Outcome,
		PrevHash:  prev,
		Reason:    e.Reason,
		Seq:       e.Seq,
		Target:    e.Target,
		Timestamp: e.Timestamp,
	})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Append adds a new hash-chained entry bound to the previous row's digest.
func (l *Log) Append(actor, action, target, outcome, reason string) Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	prev := "genesis"
	if len(l.entries) > 0 {
		prev = l.entries[len(l.entries)-1].RowHash
	}
	nowFn := l.Now
	if nowFn == nil {
		nowFn = func() time.Time { return time.Now().UTC() }
	}
	e := Entry{
		Seq:       len(l.entries) + 1,
		Timestamp: nowFn().UTC().Format(time.RFC3339),
		Actor:     actor,
		Action:    action,
		Target:    target,
		Outcome:   outcome,
		Reason:    reason,
		PrevHash:  prev,
	}
	h, err := computeRowHash(prev, e)
	if err != nil {
		panic(fmt.Sprintf("audit row hash marshal failed: %v", err))
	}
	e.RowHash = h
	l.entries = append(l.entries, e)
	return e
}

// Snapshot returns a thread-safe copy of all audit entries.
func (l *Log) Snapshot() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// TamperRow mutates a stored entry in place to simulate an insider editing history
// in tests and demos while holding the mutex cleanly.
func (l *Log) TamperRow(index int, mutate func(*Entry)) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index >= len(l.entries) {
		return fmt.Errorf("tamper index %d out of range [0,%d)", index, len(l.entries))
	}
	mutate(&l.entries[index])
	return nil
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

func (l *Log) verifyLocked() (ok bool, brokenSeq int, brokenField string) {
	if len(l.entries) == 0 {
		return false, 0, "empty_log"
	}
	prev := "genesis"
	for i, e := range l.entries {
		expectedSeq := i + 1
		if e.Seq != expectedSeq {
			return false, expectedSeq, "seq"
		}
		if e.PrevHash != prev {
			return false, e.Seq, "prev_hash"
		}
		h, err := computeRowHash(prev, e)
		if err != nil || h != e.RowHash {
			return false, e.Seq, "row_hash"
		}
		prev = e.RowHash
	}
	return true, 0, ""
}

// Verify walks the entire chain from genesis and returns the first broken sequence/field if tampered.
func (l *Log) Verify() (ok bool, brokenSeq int, brokenField string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.verifyLocked()
}

func (l *Log) exportJSONLLocked() ([]byte, error) {
	if len(l.entries) == 0 {
		return []byte{}, nil
	}
	lines := make([]string, 0, len(l.entries))
	for _, e := range l.entries {
		b, err := json.Marshal(e)
		if err != nil {
			return nil, fmt.Errorf("marshal audit entry %d: %w", e.Seq, err)
		}
		lines = append(lines, string(b))
	}
	return []byte(strings.Join(lines, "\n") + "\n"), nil
}

// ExportJSONL returns the canonical newline-terminated JSONL bytes whose SHA-256
// matches Manifest.PayloadSHA256 identically (verifiable via `shasum -a 256`).
func (l *Log) ExportJSONL() ([]byte, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.exportJSONLLocked()
}

// ExportManifest atomically verifies the chain and hashes the exact canonical JSONL bytes
// under one mutex hold (preventing TOCTOU races between verification and export).
func (l *Log) ExportManifest() (Manifest, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ok, seq, field := l.verifyLocked()
	payload, err := l.exportJSONLLocked()
	if err != nil {
		return Manifest{}, err
	}
	sum := sha256.Sum256(payload)
	nowFn := l.Now
	if nowFn == nil {
		nowFn = func() time.Time { return time.Now().UTC() }
	}
	return Manifest{
		Count:         len(l.entries),
		ChainVerified: ok,
		BrokenAtSeq:   seq,
		BrokenField:   field,
		PayloadSHA256: hex.EncodeToString(sum[:]),
		ExportedAt:    nowFn().UTC().Format(time.RFC3339),
	}, nil
}

// SignedManifest wraps Manifest + TimestampAnchor with an Ed25519 signature so auditors
// can prove neither the manifest nor the anchor was altered — even when ChainVerified is false.
type SignedManifest struct {
	Manifest        Manifest `json:"manifest"`
	TimestampAnchor string   `json:"timestamp_anchor"`
	PublicKeyHex    string   `json:"public_key_hex"`
	SignatureHex    string   `json:"signature_hex"`
}

type signedPayload struct {
	Manifest        Manifest `json:"manifest"`
	TimestampAnchor string   `json:"timestamp_anchor"`
}

// ExportSignedManifest signs both Manifest and TimestampAnchor using Ed25519.
func (l *Log) ExportSignedManifest(priv ed25519.PrivateKey, anchor string) (SignedManifest, error) {
	if len(priv) != ed25519.PrivateKeySize {
		return SignedManifest{}, errors.New("invalid ed25519 private key length")
	}
	m, err := l.ExportManifest()
	if err != nil {
		return SignedManifest{}, err
	}
	canonical, err := json.Marshal(signedPayload{
		Manifest:        m,
		TimestampAnchor: anchor,
	})
	if err != nil {
		return SignedManifest{}, fmt.Errorf("marshal signed payload: %w", err)
	}
	sig := ed25519.Sign(priv, canonical)
	pub := priv.Public().(ed25519.PublicKey)
	return SignedManifest{
		Manifest:        m,
		TimestampAnchor: anchor,
		PublicKeyHex:    hex.EncodeToString(pub),
		SignatureHex:    hex.EncodeToString(sig),
	}, nil
}

// VerifySignedManifest verifies the Ed25519 signature against a CALLER-PINNED trusted public key
// (preventing self-signed key-substitution attacks) and authenticates both Manifest and TimestampAnchor.
func VerifySignedManifest(sm SignedManifest, trustedPub ed25519.PublicKey) bool {
	if len(trustedPub) != ed25519.PublicKeySize {
		return false
	}
	embeddedPub, err := hex.DecodeString(sm.PublicKeyHex)
	if err != nil || len(embeddedPub) != ed25519.PublicKeySize {
		return false
	}
	if subtle.ConstantTimeCompare(embeddedPub, trustedPub) != 1 {
		return false
	}
	sig, err := hex.DecodeString(sm.SignatureHex)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return false
	}
	canonical, err := json.Marshal(signedPayload{
		Manifest:        sm.Manifest,
		TimestampAnchor: sm.TimestampAnchor,
	})
	if err != nil {
		return false
	}
	return ed25519.Verify(trustedPub, canonical, sig)
}
