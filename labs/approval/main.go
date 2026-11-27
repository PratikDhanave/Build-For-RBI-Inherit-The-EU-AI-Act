// Lab 01 — Five Agents, One Bad Approval (RBI FREE-AI Rec 17).
//
// Demonstrates silent governance drift: `var_calculator` was approved at medium
// risk, then refactored to declare high risk. Nobody rejected or withdrew it —
// the approval simply stopped describing the system it had approved.
//
// Run:
//
//	go run labs/approval/main.go
//	go run labs/approval/main.go --test
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type Record struct {
	Status  string
	Risk    string
	Expires *time.Time
}

var (
	now            = time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	knownRiskTiers = map[string]bool{"low": true, "medium": true, "high": true}
	in60Days       = now.AddDate(0, 0, 60)
	yesterday      = now.AddDate(0, 0, -1)

	register = map[string]*Record{
		"payment_orchestrator": {Status: "approved", Risk: "high", Expires: &in60Days},
		"kyc_orchestrator":     {Status: "approved", Risk: "high", Expires: &yesterday},
		"aml_monitor":          {Status: "rejected", Risk: "high"},
		"var_calculator":       {Status: "approved", Risk: "medium"}, // signed off at medium
		"educator":             nil,                                  // never reviewed
	}
)

func allow(agentID, declaresNow string) (bool, string) {
	risk := strings.ToLower(strings.TrimSpace(declaresNow))
	if !knownRiskTiers[risk] {
		return false, fmt.Sprintf("ErrUnknownRisk (%q)", declaresNow)
	}
	if risk != "high" {
		return true, "low/medium — no approval required"
	}
	rec := register[agentID]
	if rec == nil {
		return false, "ErrNoApproval"
	}
	if strings.ToLower(strings.TrimSpace(rec.Status)) != "approved" {
		return false, fmt.Sprintf("ErrNotApproved (%s)", rec.Status)
	}
	if rec.Expires != nil && now.After(*rec.Expires) {
		return false, "ErrExpired"
	}
	if strings.ToLower(strings.TrimSpace(rec.Risk)) != risk {
		return false, fmt.Sprintf("ErrRiskChanged (approved as %s, now declares %s)", rec.Risk, risk)
	}
	return true, "live"
}

func main() {
	today := []struct{ id, risk string }{
		{"payment_orchestrator", "high"},
		{"kyc_orchestrator", "high"},
		{"aml_monitor", "high"},
		{"var_calculator", "high"},
		{"educator", "low"},
	}

	type outcome struct {
		ok  bool
		why string
	}
	outcomes := make(map[string]outcome, len(today))

	for _, a := range today {
		ok, why := allow(a.id, a.risk)
		outcomes[a.id] = outcome{ok, why}
		status := "STOP"
		if ok {
			status = "RUN "
		}
		fmt.Printf("%s  %-22s %s\n", status, a.id, why)
	}

	fmt.Println()
	fmt.Println("var_calculator is the one I would have missed.")

	if slices.Contains(os.Args[1:], "--test") {
		must := func(cond bool, msg string) {
			if !cond {
				fmt.Fprintln(os.Stderr, "FAIL:", msg)
				os.Exit(1)
			}
		}
		must(outcomes["payment_orchestrator"].ok, "payment_orchestrator should run")
		must(!outcomes["kyc_orchestrator"].ok && outcomes["kyc_orchestrator"].why == "ErrExpired", "kyc_orchestrator should be ErrExpired")
		must(!outcomes["aml_monitor"].ok && strings.Contains(outcomes["aml_monitor"].why, "ErrNotApproved"), "aml_monitor should be ErrNotApproved")
		must(!outcomes["var_calculator"].ok && strings.Contains(outcomes["var_calculator"].why, "ErrRiskChanged"), "var_calculator should be ErrRiskChanged")
		must(outcomes["educator"].ok, "educator should run")

		// Security regression checks: padded uppercase and unknown risk tiers fail closed.
		ok, why := allow("var_calculator", " HIGH ")
		must(!ok && strings.Contains(why, "ErrRiskChanged"), "padded HIGH must catch ErrRiskChanged")
		ok, why = allow("educator", "critical")
		must(!ok && strings.Contains(why, "ErrUnknownRisk"), "unknown tier must fail with ErrUnknownRisk")
		fmt.Println("\n[OK] All approval-drift assertions passed.")
	}
}
