package incident

import (
	"encoding/json"
	"fmt"
	"time"
)

// AnnexureVIReport represents an RBI FREE-AI Annexure VI prescribed AI-incident report
// that simultaneously satisfies EU AI Act Article 73 serious-incident reporting.
type AnnexureVIReport struct {
	IncidentID          string   `json:"incident_id"`
	ReportedAt          string   `json:"reported_at"`
	RegulatedEntity     string   `json:"regulated_entity"`
	InstrumentMappings  []string `json:"instrument_mappings"`
	AgentID             string   `json:"agent_id"`
	RiskTier            string   `json:"risk_tier"`
	Category            string   `json:"category"`
	Severity            string   `json:"severity"`
	Summary             string   `json:"summary"`
	ControlTriggered    string   `json:"control_triggered"`
	CustomerImpact      string   `json:"customer_impact"`
	ContainmentAction   string   `json:"containment_action"`
	AuditChainVerified  bool     `json:"audit_chain_verified"`
	AuditManifestSHA256 string   `json:"audit_manifest_sha256"`
}

// NewAnnexureVI constructs a structured incident report ready for RBI & EU AI Act evidence packs.
func NewAnnexureVI(id, entity, agentID, riskTier, category, severity, summary, control, containment string, chainOK bool, sha string, now time.Time) AnnexureVIReport {
	return AnnexureVIReport{
		IncidentID:          id,
		ReportedAt:          now.UTC().Format(time.RFC3339),
		RegulatedEntity:     entity,
		InstrumentMappings:  []string{"RBI FREE-AI Rec 22 (Annexure VI)", "EU AI Act Article 73"},
		AgentID:             agentID,
		RiskTier:            riskTier,
		Category:            category,
		Severity:            severity,
		Summary:             summary,
		ControlTriggered:    control,
		CustomerImpact:      "None — blocked pre-dispatch by Composite Policy Gate",
		ContainmentAction:   containment,
		AuditChainVerified:  chainOK,
		AuditManifestSHA256: sha,
	}
}

// ToJSON marshals the report as pretty-printed JSON.
func (r AnnexureVIReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// ToMarkdown formats the incident report as a concise human-readable compliance brief.
func (r AnnexureVIReport) ToMarkdown() string {
	return fmt.Sprintf(
		"# RBI Annexure VI AI Incident Report — %s\n\n"+
			"- **Reported At:** %s\n"+
			"- **Regulated Entity:** %s\n"+
			"- **Agent / Model ID:** `%s` (Risk Tier: `%s`)\n"+
			"- **Category / Severity:** %s / %s\n"+
			"- **Control Triggered:** %s\n"+
			"- **Summary:** %s\n"+
			"- **Containment:** %s\n"+
			"- **Audit Chain Verified:** `%v` (`%s`)\n",
		r.IncidentID, r.ReportedAt, r.RegulatedEntity,
		r.AgentID, r.RiskTier, r.Category, r.Severity,
		r.ControlTriggered, r.Summary, r.ContainmentAction,
		r.AuditChainVerified, r.AuditManifestSHA256,
	)
}
