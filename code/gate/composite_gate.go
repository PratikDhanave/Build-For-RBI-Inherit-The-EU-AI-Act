package gate

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Message represents a trusted envelope evaluated by the Composite Policy Gate
// before dispatch to any registered agent.
//
// Trust boundary: ContainsPII, IsAdvice, Region, SubjectID, and DeclaredRisk
// MUST be stamped by deterministic gateway/DLP middleware and the live registry —
// never populated from untrusted LLM output or raw client claims.
type Message struct {
	ActorID      string // Calling user, advisor, or operator
	ActorRole    string // Normalized RBAC role (e.g., "user", "advisor", "operator")
	SubjectID    string // Customer / data subject for PII consent checks (defaults to ActorID if omitted)
	AgentID      string // Target registered agent ID
	Region       string // Target inference/processing region (e.g., "in", "adgm", "difc", "sama")
	ContainsPII  bool
	IsAdvice     bool
	Rationale    string
	ConsentLive  bool
	DeclaredRisk string // Must be "low", "medium", or "high"
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

// Allow runs every configured policy in sequence and fails closed on the first denial
// OR if the audit sink fails to persist an ALLOW decision (RBI Rec 22 / EU Art 12).
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
		if err := g.Audit.RecordAllow(msg); err != nil {
			return fmt.Errorf("DENY [audit]: failed recording allow decision: %w", err)
		}
	}
	return nil
}

// RegistryPolicy blocks unregistered shadow agents pre-dispatch (RBI Rec 23 / EU Annex IV).
type RegistryPolicy struct {
	IsRegistered func(agentID string) bool
}

func (RegistryPolicy) Name() string { return "registry" }

func (p RegistryPolicy) Check(_ context.Context, msg Message) error {
	agentID := strings.TrimSpace(msg.AgentID)
	if agentID == "" {
		return errors.New("empty target agent id")
	}
	if p.IsRegistered == nil || !p.IsRegistered(agentID) {
		return fmt.Errorf("unregistered shadow agent %q", agentID)
	}
	return nil
}

// ResidencyPolicy enforces per-message data residency (ADGM, DIFC, SAMA, IN).
// Explicit DenyRegions always take precedence over AllowedPIIRegions.
type ResidencyPolicy struct {
	AllowedPIIRegions map[string]bool
	DenyRegions       map[string]bool
}

func (ResidencyPolicy) Name() string { return "residency" }

func (p ResidencyPolicy) Check(_ context.Context, msg Message) error {
	if !msg.ContainsPII {
		return nil
	}
	region := normalizeToken(msg.Region)
	if region == "" {
		return errors.New("pii message missing target region")
	}
	if p.DenyRegions[region] {
		return fmt.Errorf("pii explicitly denied for region %s", region)
	}
	if !p.AllowedPIIRegions[region] {
		return fmt.Errorf("pii bound for region %s", region)
	}
	return nil
}

// ExplainabilityPolicy requires a non-whitespace rationale for advice or recommendations.
type ExplainabilityPolicy struct{}

func (ExplainabilityPolicy) Name() string { return "explainability" }

func (ExplainabilityPolicy) Check(_ context.Context, msg Message) error {
	if msg.IsAdvice && strings.TrimSpace(msg.Rationale) == "" {
		return errors.New("recommendation with no rationale")
	}
	return nil
}

// ConsentPolicy checks the append-only consent ledger before PII processing.
// Always evaluates against msg.SubjectID (the customer/data subject), falling back
// to msg.ActorID only when SubjectID is empty.
type ConsentPolicy struct {
	LedgerCheck func(subjectID, category string) (live bool, reason string)
}

func (ConsentPolicy) Name() string { return "consent" }

func (p ConsentPolicy) Check(_ context.Context, msg Message) error {
	if !msg.ContainsPII {
		return nil
	}
	if p.LedgerCheck != nil {
		subject := strings.TrimSpace(msg.SubjectID)
		if subject == "" {
			subject = strings.TrimSpace(msg.ActorID)
		}
		if subject == "" {
			return errors.New("pii message missing data subject id")
		}
		live, reason := p.LedgerCheck(subject, "pii")
		if !live {
			if reason != "" {
				return fmt.Errorf("no live consent (%s)", reason)
			}
			return errors.New("no live consent for this category")
		}
		return nil
	}
	if !msg.ConsentLive {
		return errors.New("no live consent for this category")
	}
	return nil
}

// RBACPolicy restricts which actor roles may invoke high-risk agents and fails closed
// on unknown or empty risk declarations.
type RBACPolicy struct {
	HighRiskRoles map[string]bool
}

func (RBACPolicy) Name() string { return "rbac" }

func (p RBACPolicy) Check(_ context.Context, msg Message) error {
	risk := normalizeToken(msg.DeclaredRisk)
	if !isKnownRiskTier(risk) {
		return fmt.Errorf("%w: %q", ErrUnknownRisk, msg.DeclaredRisk)
	}
	role := normalizeToken(msg.ActorRole)
	if risk == "high" && !p.HighRiskRoles[role] {
		return errors.New("role may not address a high-risk agent")
	}
	return nil
}
