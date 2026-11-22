package incident

import (
	"strings"
	"testing"
	"time"
)

func TestAnnexureVIGeneration(t *testing.T) {
	now := time.Date(2026, 11, 28, 10, 15, 0, 0, time.UTC)
	rep := NewAnnexureVI(
		"INC-2026-0042",
		"Sample Regulated Bank Ltd.",
		"var_calculator",
		"high",
		"Governance Drift (ErrRiskChanged)",
		"HIGH",
		"var_calculator was approved as medium risk and refactored to declare high risk without renewed CRO sign-off.",
		"ApprovalDriftGuard (RBI Rec 17)",
		"Dispatch blocked pre-execution; re-approval workflow opened.",
		true,
		"9f3c4dc6a8ff560f326e014a",
		now,
	)

	b, err := rep.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	if !strings.Contains(string(b), "EU AI Act Article 73") {
		t.Fatalf("expected EU AI Act Article 73 mapping in JSON")
	}
	md := rep.ToMarkdown()
	if !strings.Contains(md, "INC-2026-0042") || !strings.Contains(md, "var_calculator") {
		t.Fatalf("unexpected markdown output: %s", md)
	}
}
