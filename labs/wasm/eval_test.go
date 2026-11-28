package main

import (
	"os"
	"strings"
	"testing"
)

func TestEvalApprovalMatchesStandaloneAndLiveEdits(t *testing.T) {
	raw, err := os.ReadFile("../approval/main.go")
	if err != nil {
		t.Fatalf("read labs/approval/main.go: %v", err)
	}
	out, err := EvalApproval(string(raw))
	if err != nil {
		t.Fatalf("EvalApproval default failed: %v", err)
	}
	for _, want := range []string{
		"RUN   payment_orchestrator   live",
		"STOP  kyc_orchestrator       ErrExpired",
		"STOP  aml_monitor            ErrNotApproved (rejected)",
		"STOP  var_calculator         ErrRiskChanged (approved as medium, now declares high)",
		"RUN   educator               low/medium — no approval required",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}

	// Live edit on stage: promote var_calculator approval from medium -> high.
	edited := strings.Replace(string(raw),
		`"var_calculator":       {Status: "approved", Risk: "medium"}`,
		`"var_calculator":       {Status: "approved", Risk: "high"}`,
		1,
	)
	outEdited, err := EvalApproval(edited)
	if err != nil {
		t.Fatalf("EvalApproval edited failed: %v", err)
	}
	if !strings.Contains(outEdited, "RUN   var_calculator         live") {
		t.Errorf("expected var_calculator to RUN after live edit, got:\n%s", outEdited)
	}
	if !strings.Contains(out, "var_calculator is the one I would have missed.") {
		t.Errorf("default run must name the drifted agent, got:\n%s", out)
	}
	if strings.Contains(outEdited, "would have missed") {
		t.Errorf("re-approved run caught no drift and must not claim one, got:\n%s", outEdited)
	}

	// Syntax error must be surfaced by go/parser.
	if _, err := EvalApproval("package main\nfunc main( {"); err == nil {
		t.Error("expected syntax error from go/parser")
	}
}

func TestEvalAuditMatchesStandaloneAndLiveEdits(t *testing.T) {
	raw, err := os.ReadFile("../audit/main.go")
	if err != nil {
		t.Fatalf("read labs/audit/main.go: %v", err)
	}
	out, err := EvalAudit(string(raw))
	if err != nil {
		t.Fatalf("EvalAudit default failed: %v", err)
	}
	for _, want := range []string{
		"fresh chain           OK",
		"after one edited row  BROKEN at seq 2 (row_hash)",
		`"count":3,"chain_verified":false`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}

	// Live edit on stage: comment out the tamper line -> chain stays verified!
	clean := strings.Replace(string(raw), `log[1].Actor = "advisor-99"`, `// log[1].Actor = "advisor-99"`, 1)
	outClean, err := EvalAudit(clean)
	if err != nil {
		t.Fatalf("EvalAudit untampered failed: %v", err)
	}
	if !strings.Contains(outClean, "re-check, no edits    OK") || !strings.Contains(outClean, `"chain_verified":true`) {
		t.Errorf("expected untampered chain to verify true, got:\n%s", outClean)
	}
	for _, bad := range []string{"after one edited row", "does not hide the break"} {
		if strings.Contains(outClean, bad) {
			t.Errorf("untampered run must not claim an edit or a break (%q), got:\n%s", bad, outClean)
		}
	}

	// Live edit on stage: edit row 3 instead -> detects seq 3!
	row3 := strings.Replace(string(raw), `log[1].Actor = "advisor-99"`, `log[2].Action = "wire.steal"`, 1)
	outRow3, err := EvalAudit(row3)
	if err != nil {
		t.Fatalf("EvalAudit row 3 failed: %v", err)
	}
	if !strings.Contains(outRow3, "BROKEN at seq 3 (row_hash)") {
		t.Errorf("expected seq 3 break, got:\n%s", outRow3)
	}
}

func TestEvalGateConsentAndRedTeam(t *testing.T) {
	gateRaw, err := os.ReadFile("../gate_demo/main.go")
	if err != nil {
		t.Fatalf("read gate_demo: %v", err)
	}
	gateOut, err := EvalGate(string(gateRaw))
	if err != nil {
		t.Fatalf("EvalGate failed: %v", err)
	}
	if !strings.Contains(gateOut, "VERDICT: 2 allowed, 4 blocked") {
		t.Errorf("unexpected EvalGate output:\n%s", gateOut)
	}

	consentRaw, err := os.ReadFile("../consent_ledger/main.go")
	if err != nil {
		t.Fatalf("read consent_ledger: %v", err)
	}
	consentOut, err := EvalConsent(string(consentRaw))
	if err != nil {
		t.Fatalf("EvalConsent failed: %v", err)
	}
	if !strings.Contains(consentOut, "ALLOW") || !strings.Contains(consentOut, "revoked at seq 2") {
		t.Errorf("unexpected EvalConsent output:\n%s", consentOut)
	}

	rtOut, err := EvalRedTeam("ALL")
	if err != nil {
		t.Fatalf("EvalRedTeam failed: %v", err)
	}
	if !strings.Contains(rtOut, "VERDICT: 20/20 clean sweep — 13 attacks blocked, 7 safe paths allowed.") {
		t.Errorf("unexpected EvalRedTeam output:\n%s", rtOut)
	}
}
