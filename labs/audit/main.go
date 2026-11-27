// Lab 02 — Tamper-Evident SHA-256 Hash Chain & Export Manifest
// (RBI FREE-AI Rec 22 / EU AI Act Article 12).
//
// Demonstrates why an audit trail must catch its own operator: editing row 2
// invalidates the hash chain at seq 2, and the export manifest honestly reports
// `chain_verified: false` rather than hiding the break.
//
// Run:
//
//	go run labs/audit/main.go
//	go run labs/audit/main.go --test
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Entry struct {
	Seq      int    `json:"seq"`
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	Target   string `json:"target"`
	PrevHash string `json:"prev_hash"`
	RowHash  string `json:"row_hash,omitempty"`
}

func rowHash(prev string, e Entry) string {
	e.RowHash = ""
	body, _ := json.Marshal(struct {
		Prev string `json:"prev"`
		Entry
	}{Prev: prev, Entry: e})
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func appendLog(log *[]Entry, actor, action, target string) {
	prev := "genesis"
	if n := len(*log); n > 0 {
		prev = (*log)[n-1].RowHash
	}
	e := Entry{Seq: len(*log) + 1, Actor: actor, Action: action, Target: target, PrevHash: prev}
	e.RowHash = rowHash(prev, e)
	*log = append(*log, e)
}

func verify(log []Entry) (bool, int, string) {
	if len(log) == 0 {
		return false, 0, "empty_log"
	}
	prev := "genesis"
	for i, e := range log {
		if e.Seq != i+1 {
			return false, i + 1, "seq"
		}
		if e.PrevHash != prev {
			return false, e.Seq, "prev_hash"
		}
		if rowHash(prev, e) != e.RowHash {
			return false, e.Seq, "row_hash"
		}
		prev = e.RowHash
	}
	return true, 0, ""
}

func statusText(ok bool, seq int, field string) string {
	if ok {
		return "OK"
	}
	return fmt.Sprintf("BROKEN at seq %d (%s)", seq, field)
}

func main() {
	var log []Entry
	appendLog(&log, "system", "consent.grant", "user-1")
	appendLog(&log, "advisor-7", "portfolio.read", "user-1")
	appendLog(&log, "payment_orch", "refund.issue", "txn-9931")

	freshOK, freshSeq, freshField := verify(log)
	fmt.Println("fresh chain          ", statusText(freshOK, freshSeq, freshField))

	log[1].Actor = "advisor-99" // someone edits history
	tamperedOK, tamperedSeq, tamperedField := verify(log)
	fmt.Println("after one edited row ", statusText(tamperedOK, tamperedSeq, tamperedField))

	lines := make([]string, len(log))
	for i, e := range log {
		b, _ := json.Marshal(e)
		lines[i] = string(b)
	}
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	manifest, _ := json.Marshal(struct {
		Count         int    `json:"count"`
		ChainVerified bool   `json:"chain_verified"`
		SHA256        string `json:"sha256"`
	}{
		Count:         len(log),
		ChainVerified: tamperedOK,
		SHA256:        hex.EncodeToString(sum[:])[:24] + "…",
	})
	fmt.Printf("\nexport manifest       %s\n\n", manifest)
	fmt.Println("The export does not hide the break. It reports it.")

	if slices.Contains(os.Args[1:], "--test") {
		must := func(cond bool, msg string) {
			if !cond {
				fmt.Fprintln(os.Stderr, "FAIL:", msg)
				os.Exit(1)
			}
		}
		must(freshOK, "fresh chain should verify OK")
		must(!tamperedOK && tamperedSeq == 2 && tamperedField == "row_hash", "edited row must fail at seq 2 (row_hash)")
		ok, seq, field := verify(nil)
		must(!ok && seq == 0 && field == "empty_log", "empty log must fail verify")
		ok, seq, field = verify([]Entry{{Seq: 99, PrevHash: "genesis", RowHash: "x"}})
		must(!ok && seq == 1 && field == "seq", "sequence tampering must fail verify")
		fmt.Println("\n[OK] All hash-chain tamper assertions passed.")
	}
}
