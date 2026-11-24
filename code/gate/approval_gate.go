package gate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNoApproval  = errors.New("no approval record")
	ErrNotApproved = errors.New("approval status is not approved")
	ErrExpired     = errors.New("approval expired")
	ErrRiskChanged = errors.New("declared risk class drifted from approved risk class")
	ErrUnknownRisk = errors.New("unknown or empty risk class")
)

// Agent exposes its identity and currently declared runtime risk tier.
type Agent interface {
	ID() string
	DeclaredRisk() string
}

// ApprovalRecord stores the governance decision captured at sign-off time.
type ApprovalRecord struct {
	Status    string
	RiskClass string
	ExpiresAt *time.Time
}

func normalizeToken(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func isKnownRiskTier(risk string) bool {
	return risk == "low" || risk == "medium" || risk == "high"
}

// CheckApproval enforces RBI FREE-AI Rec 17 at dispatch time, catching silent
// risk-class drift (ErrRiskChanged) even when no one rejected the agent.
// Unknown or empty risk tiers fail closed (ErrUnknownRisk).
func CheckApproval(register map[string]*ApprovalRecord, a Agent, now time.Time) error {
	risk := normalizeToken(a.DeclaredRisk())
	if !isKnownRiskTier(risk) {
		return fmt.Errorf("%w: %s (%q)", ErrUnknownRisk, a.ID(), a.DeclaredRisk())
	}
	if risk != "high" {
		return nil
	}
	rec, ok := register[a.ID()]
	if !ok || rec == nil {
		return fmt.Errorf("%w: %s", ErrNoApproval, a.ID())
	}
	status := normalizeToken(rec.Status)
	if status != "approved" {
		return fmt.Errorf("%w: %s (%s)", ErrNotApproved, a.ID(), rec.Status)
	}
	if rec.ExpiresAt != nil && now.After(*rec.ExpiresAt) {
		return fmt.Errorf("%w: %s", ErrExpired, a.ID())
	}
	approvedRisk := normalizeToken(rec.RiskClass)
	if approvedRisk != risk {
		return fmt.Errorf("%w: %s approved as %s, now declares %s",
			ErrRiskChanged, a.ID(), rec.RiskClass, risk)
	}
	return nil
}

type msgAgent struct {
	id   string
	risk string
}

func (a msgAgent) ID() string           { return a.id }
func (a msgAgent) DeclaredRisk() string { return a.risk }

// ApprovalDriftPolicy wires CheckApproval directly into CompositeGate so every
// high-risk dispatch verifies active sign-off and logs denials to the audit chain.
type ApprovalDriftPolicy struct {
	Register map[string]*ApprovalRecord
	Now      func() time.Time
}

func (ApprovalDriftPolicy) Name() string { return "approval_drift" }

func (p ApprovalDriftPolicy) Check(_ context.Context, msg Message) error {
	nowFn := p.Now
	if nowFn == nil {
		nowFn = func() time.Time { return time.Now().UTC() }
	}
	return CheckApproval(p.Register, msgAgent{id: msg.AgentID, risk: msg.DeclaredRisk}, nowFn())
}
