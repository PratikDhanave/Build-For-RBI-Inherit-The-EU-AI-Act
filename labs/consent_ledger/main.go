// Lab 03 — Append-Only Consent Ledger (Slide 04: "Revoke on Tuesday, fail in-flight").
//
// Demonstrates an append-only consent ledger where a customer grants PII consent
// on Monday and revokes it on Tuesday. An in-flight PII request fails closed, and
// the immutable ledger preserves the complete GRANT -> REVOKE history so auditors
// can prove why the request failed.
//
// Run:
//
//	go run labs/consent_ledger/main.go
//	go run labs/consent_ledger/main.go --test
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type ConsentEvent struct {
	Seq      int
	Day      string
	Subject  string
	Category string
	Action   string // "GRANT" | "REVOKE"
	Reason   string
}

type AppendOnlyConsentLedger struct {
	events []ConsentEvent
}

func (l *AppendOnlyConsentLedger) Grant(subject, category, day, reason string) ConsentEvent {
	ev := ConsentEvent{Seq: len(l.events) + 1, Day: day, Subject: subject, Category: category, Action: "GRANT", Reason: reason}
	l.events = append(l.events, ev)
	return ev
}

func (l *AppendOnlyConsentLedger) Revoke(subject, category, day, reason string) ConsentEvent {
	ev := ConsentEvent{Seq: len(l.events) + 1, Day: day, Subject: subject, Category: category, Action: "REVOKE", Reason: reason}
	l.events = append(l.events, ev)
	return ev
}

func (l *AppendOnlyConsentLedger) Check(subject, category string) (bool, string) {
	for i := len(l.events) - 1; i >= 0; i-- {
		ev := l.events[i]
		if ev.Subject == subject && ev.Category == category {
			if ev.Action == "GRANT" {
				return true, fmt.Sprintf("live consent (seq %d, %s)", ev.Seq, ev.Day)
			}
			return false, fmt.Sprintf("DENY [consent] revoked at seq %d (%s: %s)", ev.Seq, ev.Day, ev.Reason)
		}
	}
	return false, "DENY [consent] no consent record"
}

func (l *AppendOnlyConsentLedger) History(subject string) []ConsentEvent {
	var out []ConsentEvent
	for _, ev := range l.events {
		if ev.Subject == subject {
			out = append(out, ev)
		}
	}
	return out
}

func main() {
	var ledger AppendOnlyConsentLedger
	ledger.Grant("user-1", "pii", "Monday 09:00", "Open Banking onboarding opt-in")

	monOK, monWhy := ledger.Check("user-1", "pii")
	fmt.Printf("Monday  in-flight PII request : %s — %s\n", verdict(monOK), monWhy)

	ledger.Revoke("user-1", "pii", "Tuesday 14:30", "customer revoked on Tuesday")

	tueOK, tueWhy := ledger.Check("user-1", "pii")
	fmt.Printf("Tuesday in-flight PII request : %s — %s\n", verdict(tueOK), tueWhy)

	fmt.Println("\nAppend-only history for user-1:")
	for _, ev := range ledger.History("user-1") {
		fmt.Printf("  seq %d | %-14s | %-6s | %s\n", ev.Seq, ev.Day, ev.Action, ev.Reason)
	}

	if slices.Contains(os.Args[1:], "--test") {
		if !monOK {
			fmt.Fprintln(os.Stderr, "FAIL: Monday request should be allowed")
			os.Exit(1)
		}
		if tueOK || !strings.Contains(tueWhy, "revoked at seq 2") {
			fmt.Fprintf(os.Stderr, "FAIL: Tuesday request must fail with seq 2 revocation, got ok=%v why=%q\n", tueOK, tueWhy)
			os.Exit(1)
		}
		if got := len(ledger.History("user-1")); got != 2 {
			fmt.Fprintf(os.Stderr, "FAIL: expected 2 history events, got %d\n", got)
			os.Exit(1)
		}
		fmt.Println("\n[OK] Append-only consent ledger assertions passed.")
	}
}

func verdict(ok bool) string {
	if ok {
		return "ALLOW"
	}
	return "DENY "
}
