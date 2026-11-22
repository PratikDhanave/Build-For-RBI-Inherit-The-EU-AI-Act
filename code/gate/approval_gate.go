package gate

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNoApproval  = errors.New("no approval record")
	ErrNotApproved = errors.New("approval status is not approved")
	ErrExpired     = errors.New("approval expired")
	ErrRiskChanged = errors.New("declared risk class drifted from approved risk class")
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

// CheckApproval enforces RBI FREE-AI Rec 17 at dispatch time, catching silent
// risk-class drift (ErrRiskChanged) even when no one rejected the agent.
func CheckApproval(register map[string]*ApprovalRecord, a Agent, now time.Time) error {
	risk := a.DeclaredRisk()
	if risk != "high" {
		return nil
	}
	rec, ok := register[a.ID()]
	if !ok || rec == nil {
		return fmt.Errorf("%w: %s", ErrNoApproval, a.ID())
	}
	if rec.Status != "approved" {
		return fmt.Errorf("%w: %s (%s)", ErrNotApproved, a.ID(), rec.Status)
	}
	if rec.ExpiresAt != nil && now.After(*rec.ExpiresAt) {
		return fmt.Errorf("%w: %s", ErrExpired, a.ID())
	}
	if rec.RiskClass != risk {
		return fmt.Errorf("%w: %s approved as %s, now declares %s",
			ErrRiskChanged, a.ID(), rec.RiskClass, risk)
	}
	return nil
}
