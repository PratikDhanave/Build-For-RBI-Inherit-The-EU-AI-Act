package gate

import (
	"context"
	"errors"
	"fmt"
)

// Message represents an envelope evaluated by the Composite Policy Gate
// before dispatch to any of the 62 registered agents.
type Message struct {
	ActorID      string
	ActorRole    string
	AgentID      string
	Region       string
	ContainsPII  bool
	IsAdvice     bool
	Rationale    string
	ConsentLive  bool
	DeclaredRisk string
}

// Policy checks a single deterministic governance rule pre-dispatch.
type Policy interface {
	Name() string
	Check(ctx context.Context, msg Message) error
}

// AuditLogger records every allow/deny decision into the tamper-evident chain.
type AuditLogger interface {
	RecordAllow(msg Message) error
	RecordDeny(msg Message, policy, reason string) error
}

// CompositeGate evaluates an ordered list of policies before any agent sees msg.
type CompositeGate struct {
	Policies []Policy
	Audit    AuditLogger
}

// Allow runs every configured policy in sequence and fails closed on the first denial.
func (g *CompositeGate) Allow(ctx context.Context, msg Message) error {
	for _, p := range g.Policies {
		if err := p.Check(ctx, msg); err != nil {
			if g.Audit != nil {
				_ = g.Audit.RecordDeny(msg, p.Name(), err.Error())
			}
			return fmt.Errorf("DENY [%s]: %w", p.Name(), err)
		}
	}
	if g.Audit != nil {
		_ = g.Audit.RecordAllow(msg)
	}
	return nil
}

// ResidencyPolicy enforces per-message data residency (ADGM, DIFC, SAMA, IN).
type ResidencyPolicy struct {
	AllowedPIIRegions map[string]bool
}

func (p ResidencyPolicy) Name() string { return "residency" }

func (p ResidencyPolicy) Check(_ context.Context, msg Message) error {
	if msg.ContainsPII && !p.AllowedPIIRegions[msg.Region] {
		return fmt.Errorf("pii bound for region %s", msg.Region)
	}
	return nil
}

// ExplainabilityPolicy requires a rationale for advice or recommendations.
type ExplainabilityPolicy struct{}

func (ExplainabilityPolicy) Name() string { return "explainability" }

func (ExplainabilityPolicy) Check(_ context.Context, msg Message) error {
	if msg.IsAdvice && msg.Rationale == "" {
		return errors.New("recommendation with no rationale")
	}
	return nil
}

// ConsentPolicy checks the append-only consent ledger before PII processing.
type ConsentPolicy struct{}

func (ConsentPolicy) Name() string { return "consent" }

func (ConsentPolicy) Check(_ context.Context, msg Message) error {
	if msg.ContainsPII && !msg.ConsentLive {
		return errors.New("no live consent for this category")
	}
	return nil
}

// RBACPolicy restricts which actor roles may invoke high-risk agents.
type RBACPolicy struct {
	HighRiskRoles map[string]bool
}

func (p RBACPolicy) Name() string { return "rbac" }

func (p RBACPolicy) Check(_ context.Context, msg Message) error {
	if msg.DeclaredRisk == "high" && !p.HighRiskRoles[msg.ActorRole] {
		return errors.New("role may not address a high-risk agent")
	}
	return nil
}
