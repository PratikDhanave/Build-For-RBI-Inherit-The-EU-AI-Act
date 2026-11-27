// Composite Policy Gate Demo — One Gate, Every Message, Before Dispatch.
//
// Demonstrates how a single pre-dispatch gate enforces residency, explainability,
// RBAC, and append-only consent policies across 62 registered agents without
// duplicating compliance logic inside individual agent code.
//
// Run:
//
//	go run labs/gate_demo/main.go
//	go run labs/gate_demo/main.go --test
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Message struct {
	Label        string
	ActorRole    string
	TargetRisk   string
	ContainsPII  bool
	TargetRegion string
	IsAdvice     bool
	HasRationale bool
	ConsentLive  bool
}

var (
	allowedPIIRegions = map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true}
	knownRiskTiers    = map[string]bool{"low": true, "medium": true, "high": true}
	highRiskRoles     = map[string]bool{"advisor": true, "operator": true}
)

// evaluateGate evaluates ordered policies before any agent sees the message (fail-closed).
func evaluateGate(msg Message) (bool, string) {
	risk := strings.ToLower(strings.TrimSpace(msg.TargetRisk))
	role := strings.ToLower(strings.TrimSpace(msg.ActorRole))
	region := strings.ToLower(strings.TrimSpace(msg.TargetRegion))
	if region == "" {
		region = "in"
	}

	if !knownRiskTiers[risk] {
		return false, fmt.Sprintf("[rbac] unknown risk tier %q", msg.TargetRisk)
	}

	// Policy 1: Per-message data residency
	if msg.ContainsPII && !allowedPIIRegions[region] {
		return false, fmt.Sprintf("[residency] pii bound for region %s", region)
	}

	// Policy 2: Explainability requirement for advice/recommendations
	if msg.IsAdvice && !msg.HasRationale {
		return false, "[explainability] recommendation with no rationale"
	}

	// Policy 3: Role-based access control for high-risk agents
	if risk == "high" && !highRiskRoles[role] {
		return false, "[rbac] role may not address a high-risk agent"
	}

	// Policy 4: Append-only live consent verification
	if msg.ContainsPII && !msg.ConsentLive {
		return false, "[consent] no live consent for this category"
	}

	return true, "dispatched to agent"
}

var messages = []Message{
	{Label: "ordinary question", ActorRole: "user", TargetRisk: "low", TargetRegion: "in", HasRationale: true, ConsentLive: true},
	{Label: "PAN to a US model", ActorRole: "advisor", TargetRisk: "medium", ContainsPII: true, TargetRegion: "us", HasRationale: true, ConsentLive: true},
	{Label: "advice with no reason", ActorRole: "advisor", TargetRisk: "medium", TargetRegion: "in", IsAdvice: true, HasRationale: false, ConsentLive: true},
	{Label: "user hits a high agent", ActorRole: "user", TargetRisk: "high", TargetRegion: "in", HasRationale: true, ConsentLive: true},
	{Label: "PII, consent revoked", ActorRole: "advisor", TargetRisk: "medium", ContainsPII: true, TargetRegion: "in", HasRationale: true, ConsentLive: false},
	{Label: "balance lookup", ActorRole: "user", TargetRisk: "low", TargetRegion: "in", HasRationale: true, ConsentLive: true},
}

func main() {
	results := make(map[string]bool, len(messages))
	for _, msg := range messages {
		allowed, reason := evaluateGate(msg)
		results[msg.Label] = allowed
		verdict := "DENY "
		detail := " " + reason
		if allowed {
			verdict = "ALLOW"
			detail = ""
		}
		fmt.Printf("%-26s %s%s\n", msg.Label, verdict, detail)
	}

	fmt.Println()
	fmt.Println("Nothing above is agent code. The gate runs before any agent sees it.")

	if slices.Contains(os.Args[1:], "--test") {
		expected := map[string]bool{
			"ordinary question":      true,
			"PAN to a US model":      false,
			"advice with no reason":  false,
			"user hits a high agent": false,
			"PII, consent revoked":   false,
			"balance lookup":         true,
		}
		for label, want := range expected {
			if got := results[label]; got != want {
				fmt.Fprintf(os.Stderr, "FAIL: %q allowed=%v, want %v\n", label, got, want)
				os.Exit(1)
			}
		}
		// Security regression checks: uppercase/padded and unknown risk tiers fail closed.
		if ok, _ := evaluateGate(Message{Label: "upper", ActorRole: "user", TargetRisk: " HIGH ", TargetRegion: "in", ConsentLive: true}); ok {
			fmt.Fprintln(os.Stderr, "FAIL: padded uppercase HIGH risk must fail closed")
			os.Exit(1)
		}
		if ok, _ := evaluateGate(Message{Label: "unknown", ActorRole: "advisor", TargetRisk: "critical", TargetRegion: "in", ConsentLive: true}); ok {
			fmt.Fprintln(os.Stderr, "FAIL: unknown risk tier must fail closed")
			os.Exit(1)
		}
		fmt.Println("\n[OK] All composite policy gate assertions passed.")
	}
}
