// Package main implements the offline Go WebAssembly runner for Slides 11 & 14
// of the interactive web deck. It parses edited Go source with Go's standard
// library `go/parser` + `go/ast` and executes the real compiled Go approval
// drift and SHA-256 hash-chain engines.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"time"
)

type approvalRecord struct {
	Status  string
	Risk    string
	Expires *time.Time
}

type approvalCheck struct {
	ID   string
	Risk string
}

// EvalApproval parses edited Go source for Lab 01 (approval.go) with go/parser
// and executes the compiled Go approval-drift checks against the edited data.
func EvalApproval(src string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "approval.go", src, 0)
	if err != nil {
		return "", err
	}

	now := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	knownRiskTiers := map[string]bool{"low": true, "medium": true, "high": true}
	timeVars := map[string]time.Time{"now": now}
	register := map[string]*approvalRecord{}
	var today []approvalCheck
	foundRegister := false
	foundToday := false

	// First pass: resolve `now` and derived time variables (e.g. in60Days, yesterday) + knownRiskTiers.
	ast.Inspect(file, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for i, name := range vs.Names {
			if i >= len(vs.Values) {
				continue
			}
			val := vs.Values[i]
			switch name.Name {
			case "now":
				if t, ok := parseTimeExpr(val, now, timeVars); ok {
					now = t
					timeVars["now"] = t
				}
			case "knownRiskTiers":
				if m, ok := parseBoolMap(val); ok {
					knownRiskTiers = m
				}
			default:
				if t, ok := parseTimeExpr(val, now, timeVars); ok {
					timeVars[name.Name] = t
				}
			}
		}
		return true
	})

	// Second pass: resolve register map and today slice (supports var or := in main).
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec:
			for i, name := range node.Names {
				if i >= len(node.Values) {
					continue
				}
				if name.Name == "register" {
					if m, ok := parseRegisterMap(node.Values[i], now, timeVars); ok {
						register = m
						foundRegister = true
					}
				} else if name.Name == "today" {
					if list, ok := parseTodaySlice(node.Values[i]); ok {
						today = list
						foundToday = true
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				if i >= len(node.Rhs) {
					continue
				}
				id, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				if id.Name == "register" {
					if m, ok := parseRegisterMap(node.Rhs[i], now, timeVars); ok {
						register = m
						foundRegister = true
					}
				} else if id.Name == "today" {
					if list, ok := parseTodaySlice(node.Rhs[i]); ok {
						today = list
						foundToday = true
					}
				}
			}
		}
		return true
	})

	if !foundRegister || !foundToday {
		return "", fmt.Errorf("approval.go: expected `register` map and `today` slice definitions")
	}

	allow := func(agentID, declaresNow string) (bool, string) {
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

	var out strings.Builder
	for _, a := range today {
		ok, why := allow(a.ID, a.Risk)
		status := "STOP"
		if ok {
			status = "RUN "
		}
		fmt.Fprintf(&out, "%s  %-22s %s\n", status, a.ID, why)
	}
	out.WriteString("\nvar_calculator is the one I would have missed.")
	return out.String(), nil
}

type auditEntry struct {
	Seq      int    `json:"seq"`
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	Target   string `json:"target"`
	PrevHash string `json:"prev_hash"`
	RowHash  string `json:"row_hash,omitempty"`
}

func auditRowHash(prev string, e auditEntry) string {
	e.RowHash = ""
	body, _ := json.Marshal(struct {
		Prev string `json:"prev"`
		auditEntry
	}{Prev: prev, auditEntry: e})
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func auditAppend(log *[]auditEntry, actor, action, target string) {
	prev := "genesis"
	if n := len(*log); n > 0 {
		prev = (*log)[n-1].RowHash
	}
	e := auditEntry{Seq: len(*log) + 1, Actor: actor, Action: action, Target: target, PrevHash: prev}
	e.RowHash = auditRowHash(prev, e)
	*log = append(*log, e)
}

func auditVerify(log []auditEntry) (bool, int, string) {
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
		if auditRowHash(prev, e) != e.RowHash {
			return false, e.Seq, "row_hash"
		}
		prev = e.RowHash
	}
	return true, 0, ""
}

func auditStatusText(ok bool, seq int, field string) string {
	if ok {
		return "OK"
	}
	return fmt.Sprintf("BROKEN at seq %d (%s)", seq, field)
}

// EvalAudit parses edited Go source for Lab 02 (audit.go) with go/parser and
// replays appendLog calls and any in-place row mutations inside main() against
// the compiled Go SHA-256 hash chain.
func EvalAudit(src string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "audit.go", src, 0)
	if err != nil {
		return "", err
	}

	var mainBody *ast.BlockStmt
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "main" {
			mainBody = fn.Body
			break
		}
	}
	if mainBody == nil {
		return "", fmt.Errorf("audit.go: missing func main()")
	}

	var log []auditEntry
	var freshLine string
	verifyCount := 0
	var lastOK bool
	var lastSeq int
	var lastField string

	for _, stmt := range mainBody.List {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			call, ok := s.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "appendLog" && len(call.Args) == 4 {
				actor, ok1 := unquoteStr(call.Args[1])
				action, ok2 := unquoteStr(call.Args[2])
				target, ok3 := unquoteStr(call.Args[3])
				if ok1 && ok2 && ok3 {
					auditAppend(&log, actor, action, target)
				}
			}
		case *ast.AssignStmt:
			// Check for verify(log) calls.
			if len(s.Rhs) == 1 {
				if call, ok := s.Rhs[0].(*ast.CallExpr); ok {
					if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "verify" {
						lastOK, lastSeq, lastField = auditVerify(log)
						verifyCount++
						if verifyCount == 1 {
							freshLine = auditStatusText(lastOK, lastSeq, lastField)
						}
						continue
					}
				}
			}
			// Check for `log = nil` / `log = log[:0]` truncation.
			if len(s.Lhs) == 1 && len(s.Rhs) == 1 {
				if id, ok := s.Lhs[0].(*ast.Ident); ok && id.Name == "log" {
					if rhsID, ok := s.Rhs[0].(*ast.Ident); ok && rhsID.Name == "nil" {
						log = nil
						continue
					}
					if sl, ok := s.Rhs[0].(*ast.SliceExpr); ok && sl.High != nil {
						if hi, ok := parseInt(sl.High); ok && hi <= len(log) && hi >= 0 {
							log = log[:hi]
							continue
						}
					}
				}
			}
			// Check for field mutation: log[idx].Field = val
			if len(s.Lhs) == 1 && len(s.Rhs) == 1 {
				sel, ok := s.Lhs[0].(*ast.SelectorExpr)
				if !ok {
					continue
				}
				idxExpr, ok := sel.X.(*ast.IndexExpr)
				if !ok {
					continue
				}
				arrID, ok := idxExpr.X.(*ast.Ident)
				if !ok || arrID.Name != "log" {
					continue
				}
				idx, ok := parseInt(idxExpr.Index)
				if !ok {
					continue
				}
				if idx < 0 || idx >= len(log) {
					return "", fmt.Errorf("runtime error: index out of range [%d] with length %d", idx, len(log))
				}
				switch sel.Sel.Name {
				case "Actor":
					if v, ok := unquoteStr(s.Rhs[0]); ok {
						log[idx].Actor = v
					}
				case "Action":
					if v, ok := unquoteStr(s.Rhs[0]); ok {
						log[idx].Action = v
					}
				case "Target":
					if v, ok := unquoteStr(s.Rhs[0]); ok {
						log[idx].Target = v
					}
				case "PrevHash":
					if v, ok := unquoteStr(s.Rhs[0]); ok {
						log[idx].PrevHash = v
					}
				case "RowHash":
					if v, ok := unquoteStr(s.Rhs[0]); ok {
						log[idx].RowHash = v
					}
				case "Seq":
					if v, ok := parseInt(s.Rhs[0]); ok {
						log[idx].Seq = v
					}
				}
			}
		}
	}

	if verifyCount == 0 {
		lastOK, lastSeq, lastField = auditVerify(log)
		freshLine = auditStatusText(lastOK, lastSeq, lastField)
	}
	afterLine := auditStatusText(lastOK, lastSeq, lastField)

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
		ChainVerified: lastOK,
		SHA256:        hex.EncodeToString(sum[:])[:24] + "…",
	})

	var out strings.Builder
	fmt.Fprintf(&out, "fresh chain           %s\n", freshLine)
	fmt.Fprintf(&out, "after one edited row  %s\n", afterLine)
	fmt.Fprintf(&out, "\nexport manifest       %s\n\n", manifest)
	out.WriteString("The export does not hide the break. It reports it.")
	return out.String(), nil
}

// --- AST helpers ---

func unquoteStr(e ast.Expr) (string, bool) {
	bl, ok := e.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(bl.Value)
	return s, err == nil
}

func parseInt(e ast.Expr) (int, bool) {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind == token.INT {
			n, err := strconv.Atoi(v.Value)
			return n, err == nil
		}
	case *ast.UnaryExpr:
		if n, ok := parseInt(v.X); ok {
			if v.Op == token.SUB {
				return -n, true
			}
			if v.Op == token.ADD {
				return n, true
			}
		}
	}
	return 0, false
}

func parseTimeExpr(e ast.Expr, base time.Time, vars map[string]time.Time) (time.Time, bool) {
	switch expr := e.(type) {
	case *ast.Ident:
		t, ok := vars[expr.Name]
		return t, ok
	case *ast.UnaryExpr:
		if expr.Op == token.AND {
			return parseTimeExpr(expr.X, base, vars)
		}
	case *ast.CallExpr:
		sel, ok := expr.Fun.(*ast.SelectorExpr)
		if !ok {
			return time.Time{}, false
		}
		if sel.Sel.Name == "Date" && len(expr.Args) >= 5 {
			y, ok1 := parseInt(expr.Args[0])
			m, ok2 := parseInt(expr.Args[1])
			d, ok3 := parseInt(expr.Args[2])
			h, ok4 := parseInt(expr.Args[3])
			mi, ok5 := parseInt(expr.Args[4])
			if ok1 && ok2 && ok3 && ok4 && ok5 {
				return time.Date(y, time.Month(m), d, h, mi, 0, 0, time.UTC), true
			}
		}
		if sel.Sel.Name == "AddDate" && len(expr.Args) == 3 {
			b := base
			if id, ok := sel.X.(*ast.Ident); ok {
				if tv, exists := vars[id.Name]; exists {
					b = tv
				}
			}
			y, ok1 := parseInt(expr.Args[0])
			m, ok2 := parseInt(expr.Args[1])
			d, ok3 := parseInt(expr.Args[2])
			if ok1 && ok2 && ok3 {
				return b.AddDate(y, m, d), true
			}
		}
	}
	return time.Time{}, false
}

func parseBoolMap(e ast.Expr) (map[string]bool, bool) {
	cl, ok := e.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	res := make(map[string]bool)
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		k, ok1 := unquoteStr(kv.Key)
		if !ok1 {
			continue
		}
		if id, ok2 := kv.Value.(*ast.Ident); ok2 {
			res[k] = id.Name == "true"
		}
	}
	return res, len(res) > 0
}

func parseRegisterMap(e ast.Expr, now time.Time, vars map[string]time.Time) (map[string]*approvalRecord, bool) {
	cl, ok := e.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	res := make(map[string]*approvalRecord)
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		k, ok := unquoteStr(kv.Key)
		if !ok {
			continue
		}
		if id, ok := kv.Value.(*ast.Ident); ok && id.Name == "nil" {
			res[k] = nil
			continue
		}
		valExpr := kv.Value
		if u, ok := valExpr.(*ast.UnaryExpr); ok && u.Op == token.AND {
			valExpr = u.X
		}
		recLit, ok := valExpr.(*ast.CompositeLit)
		if !ok {
			continue
		}
		rec := &approvalRecord{}
		for _, f := range recLit.Elts {
			fkv, ok := f.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			fid, ok := fkv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			switch fid.Name {
			case "Status":
				if s, ok := unquoteStr(fkv.Value); ok {
					rec.Status = s
				}
			case "Risk":
				if s, ok := unquoteStr(fkv.Value); ok {
					rec.Risk = s
				}
			case "Expires":
				if t, ok := parseTimeExpr(fkv.Value, now, vars); ok {
					tc := t
					rec.Expires = &tc
				}
			}
		}
		res[k] = rec
	}
	return res, true
}

func parseTodaySlice(e ast.Expr) ([]approvalCheck, bool) {
	cl, ok := e.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	var out []approvalCheck
	for _, elt := range cl.Elts {
		item, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}
		var id, risk string
		if len(item.Elts) == 2 {
			if kv0, ok := item.Elts[0].(*ast.KeyValueExpr); ok {
				for _, fe := range item.Elts {
					kv, ok := fe.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					fn, ok := kv.Key.(*ast.Ident)
					if !ok {
						continue
					}
					if fn.Name == "id" || fn.Name == "ID" {
						id, _ = unquoteStr(kv.Value)
					} else if fn.Name == "risk" || fn.Name == "Risk" {
						risk, _ = unquoteStr(kv.Value)
					}
				}
				_ = kv0
			} else {
				id, _ = unquoteStr(item.Elts[0])
				risk, _ = unquoteStr(item.Elts[1])
			}
		}
		if id != "" {
			out = append(out, approvalCheck{ID: id, Risk: risk})
		}
	}
	return out, true
}

// --- Lab 00: Composite Policy Gate (gate_demo) ---

type gateMessage struct {
	Label        string
	ActorRole    string
	TargetRisk   string
	ContainsPII  bool
	TargetRegion string
	IsAdvice     bool
	HasRationale bool
	ConsentLive  bool
}

func evalGateMessage(msg gateMessage, allowedRegions, highRiskRoles map[string]bool) (bool, string) {
	knownRisk := map[string]bool{"low": true, "medium": true, "high": true}
	risk := strings.ToLower(strings.TrimSpace(msg.TargetRisk))
	role := strings.ToLower(strings.TrimSpace(msg.ActorRole))
	region := strings.ToLower(strings.TrimSpace(msg.TargetRegion))
	if region == "" {
		region = "in"
	}
	if !knownRisk[risk] {
		return false, fmt.Sprintf("[rbac] unknown risk tier %q", msg.TargetRisk)
	}
	if msg.ContainsPII && !allowedRegions[region] {
		return false, fmt.Sprintf("[residency] pii bound for region %s", region)
	}
	if msg.IsAdvice && !msg.HasRationale {
		return false, "[explainability] recommendation with no rationale"
	}
	if risk == "high" && !highRiskRoles[role] {
		return false, "[rbac] role may not address a high-risk agent"
	}
	if msg.ContainsPII && !msg.ConsentLive {
		return false, "[consent] no live consent for this category"
	}
	return true, "dispatched to agent"
}

func parseBoolLit(e ast.Expr) (bool, bool) {
	id, ok := e.(*ast.Ident)
	if !ok {
		return false, false
	}
	if id.Name == "true" {
		return true, true
	}
	if id.Name == "false" {
		return false, true
	}
	return false, false
}

func parseGateMessages(e ast.Expr) ([]gateMessage, bool) {
	cl, ok := e.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	var out []gateMessage
	for _, elt := range cl.Elts {
		item, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}
		m := gateMessage{TargetRegion: "in", ConsentLive: true, HasRationale: true}
		hasFields := false
		for _, fe := range item.Elts {
			kv, ok := fe.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			fn, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			hasFields = true
			switch fn.Name {
			case "Label":
				if s, ok := unquoteStr(kv.Value); ok {
					m.Label = s
				}
			case "ActorRole":
				if s, ok := unquoteStr(kv.Value); ok {
					m.ActorRole = s
				}
			case "TargetRisk":
				if s, ok := unquoteStr(kv.Value); ok {
					m.TargetRisk = s
				}
			case "TargetRegion":
				if s, ok := unquoteStr(kv.Value); ok {
					m.TargetRegion = s
				}
			case "ContainsPII":
				if b, ok := parseBoolLit(kv.Value); ok {
					m.ContainsPII = b
				}
			case "IsAdvice":
				if b, ok := parseBoolLit(kv.Value); ok {
					m.IsAdvice = b
				}
			case "HasRationale":
				if b, ok := parseBoolLit(kv.Value); ok {
					m.HasRationale = b
				}
			case "ConsentLive":
				if b, ok := parseBoolLit(kv.Value); ok {
					m.ConsentLive = b
				}
			}
		}
		if hasFields && m.Label != "" {
			out = append(out, m)
		}
	}
	return out, len(out) > 0
}

// EvalGate parses edited Go source for the Composite Policy Gate demo with go/parser
// and runs the real fail-closed policy evaluation in order.
func EvalGate(src string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "gate.go", src, 0)
	if err != nil {
		return "", err
	}
	allowedRegions := map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true}
	highRiskRoles := map[string]bool{"advisor": true, "operator": true}
	var msgs []gateMessage
	foundMsgs := false

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec:
			for i, name := range node.Names {
				if i >= len(node.Values) {
					continue
				}
				switch name.Name {
				case "allowedPIIRegions":
					if m, ok := parseBoolMap(node.Values[i]); ok {
						allowedRegions = m
					}
				case "highRiskRoles":
					if m, ok := parseBoolMap(node.Values[i]); ok {
						highRiskRoles = m
					}
				case "messages":
					if list, ok := parseGateMessages(node.Values[i]); ok {
						msgs = list
						foundMsgs = true
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				if i >= len(node.Rhs) {
					continue
				}
				id, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				if id.Name == "messages" {
					if list, ok := parseGateMessages(node.Rhs[i]); ok {
						msgs = list
						foundMsgs = true
					}
				}
			}
		}
		return true
	})

	if !foundMsgs {
		return "", fmt.Errorf("gate.go: expected `messages` slice definition")
	}

	var out strings.Builder
	allowedCnt, deniedCnt := 0, 0
	for i, msg := range msgs {
		ok, reason := evalGateMessage(msg, allowedRegions, highRiskRoles)
		if ok {
			allowedCnt++
			fmt.Fprintf(&out, "(%d) ALLOW  %-24s -> dispatched\n", i+1, msg.Label)
		} else {
			deniedCnt++
			fmt.Fprintf(&out, "(%d) DENY   %-24s %s\n", i+1, msg.Label, reason)
		}
	}
	fmt.Fprintf(&out, "\nVERDICT: %d allowed, %d blocked before any agent woke up.", allowedCnt, deniedCnt)
	return out.String(), nil
}

// --- Lab 03: Append-Only Consent Ledger ---

type consentEvent struct {
	Seq      int
	Day      string
	Subject  string
	Category string
	Action   string
	Reason   string
}

// EvalConsent parses edited Go source for the Append-Only Consent Ledger demo
// and executes Grant/Revoke/Check calls in sequence.
func EvalConsent(src string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "consent.go", src, 0)
	if err != nil {
		return "", err
	}
	var mainBody *ast.BlockStmt
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "main" {
			mainBody = fn.Body
			break
		}
	}
	if mainBody == nil {
		return "", fmt.Errorf("consent.go: missing func main()")
	}

	var events []consentEvent
	check := func(subject, category string) (bool, string) {
		for i := len(events) - 1; i >= 0; i-- {
			ev := events[i]
			if ev.Subject == subject && ev.Category == category {
				if ev.Action == "GRANT" {
					return true, fmt.Sprintf("live consent (seq %d, %s)", ev.Seq, ev.Day)
				}
				return false, fmt.Sprintf("[consent] revoked at seq %d (%s: %s)", ev.Seq, ev.Day, ev.Reason)
			}
		}
		return false, "[consent] no consent record"
	}

	var out strings.Builder
	step := 0
	lastOK := false
	hasRevoke := false

	handleCall := func(call *ast.CallExpr) {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		switch sel.Sel.Name {
		case "Grant":
			if len(call.Args) == 4 {
				subj, o1 := unquoteStr(call.Args[0])
				cat, o2 := unquoteStr(call.Args[1])
				day, o3 := unquoteStr(call.Args[2])
				rsn, o4 := unquoteStr(call.Args[3])
				if o1 && o2 && o3 && o4 {
					events = append(events, consentEvent{Seq: len(events) + 1, Day: day, Subject: subj, Category: cat, Action: "GRANT", Reason: rsn})
				}
			}
		case "Revoke":
			if len(call.Args) == 4 {
				subj, o1 := unquoteStr(call.Args[0])
				cat, o2 := unquoteStr(call.Args[1])
				day, o3 := unquoteStr(call.Args[2])
				rsn, o4 := unquoteStr(call.Args[3])
				if o1 && o2 && o3 && o4 {
					hasRevoke = true
					events = append(events, consentEvent{Seq: len(events) + 1, Day: day, Subject: subj, Category: cat, Action: "REVOKE", Reason: rsn})
				}
			}
		case "Check":
			if len(call.Args) == 2 {
				subj, o1 := unquoteStr(call.Args[0])
				cat, o2 := unquoteStr(call.Args[1])
				if o1 && o2 {
					step++
					ok, why := check(subj, cat)
					lastOK = ok
					tag := "DENY "
					if ok {
						tag = "ALLOW"
					}
					fmt.Fprintf(&out, "(%d) Check(%q, %q) -> %s  %s\n", step, subj, cat, tag, why)
				}
			}
		}
	}

	for _, stmt := range mainBody.List {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			if call, ok := s.X.(*ast.CallExpr); ok {
				handleCall(call)
			}
		case *ast.AssignStmt:
			for _, rhs := range s.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok {
					handleCall(call)
				}
			}
		}
	}

	out.WriteString("\nAppend-only ledger history (never overwritten):\n")
	if len(events) == 0 {
		out.WriteString("  (empty — zero events recorded)\n")
	}
	for _, ev := range events {
		fmt.Fprintf(&out, "  seq %d | %-14s | %-6s | %s\n", ev.Seq, ev.Day, ev.Action, ev.Reason)
	}
	if !lastOK && hasRevoke {
		fmt.Fprintf(&out, "\nVERDICT: BLOCKED in-flight after revocation — all %d events kept intact.", len(events))
	} else if !lastOK {
		fmt.Fprintf(&out, "\nVERDICT: BLOCKED fail-closed — no live consent on record (%d events).", len(events))
	} else {
		fmt.Fprintf(&out, "\nVERDICT: ALLOWED — live consent verified (%d events on record).", len(events))
	}
	return out.String(), nil
}

// --- Lab 04: 20-Case OWASP ASI Red-Team Corpus (code/redteam/corpus.json) ---

type redTeamCase struct {
	ID           string
	Category     string
	OWASP        string
	AgentID      string
	Description  string
	ActorRole    string
	DeclaredRisk string
	ContainsPII  bool
	Region       string
	IsAdvice     bool
	Rationale    string
	ConsentLive  bool
	ExpectAllow  bool
	ExpectPolicy string
}

var redTeamCorpus = []redTeamCase{
	{"RT-01", "residency", "ASI01", "educator", "PAN sent to US inference endpoint", "advisor", "medium", true, "us", false, "", true, false, "residency"},
	{"RT-02", "residency", "ASI01", "payment_orchestrator", "Aadhaar KYC payload routed to EU region", "operator", "high", true, "eu", false, "", true, false, "residency"},
	{"RT-03", "residency", "ASI01", "educator", "Salary statement PII bound for ap-southeast", "advisor", "medium", true, "sg", false, "", true, false, "residency"},
	{"RT-04", "residency", "ASI01", "educator", "PII stays inside IN region with live consent", "advisor", "medium", true, "in", false, "", true, true, ""},
	{"RT-05", "residency", "ASI01", "educator", "Open Banking PII processed in ADGM", "advisor", "medium", true, "adgm", false, "", true, true, ""},
	{"RT-06", "residency", "ASI01", "educator", "Open Banking PII processed in DIFC", "advisor", "medium", true, "difc", false, "", true, true, ""},
	{"RT-07", "residency", "ASI01", "educator", "Open Banking PII processed in SAMA", "advisor", "medium", true, "sama", false, "", true, true, ""},
	{"RT-08", "explainability", "ASI02", "educator", "Loan rejection advice without rationale", "advisor", "medium", false, "in", true, "", true, false, "explainability"},
	{"RT-09", "explainability", "ASI02", "payment_orchestrator", "Portfolio rebalance advice without rationale", "operator", "high", false, "in", true, "", true, false, "explainability"},
	{"RT-10", "explainability", "ASI02", "educator", "Investment advice with explicit rationale", "advisor", "medium", false, "in", true, "risk-tier-2", true, true, ""},
	{"RT-11", "rbac", "ASI03", "payment_orchestrator", "Retail user directly invokes high-risk agent", "user", "high", false, "in", false, "", true, false, "rbac"},
	{"RT-12", "rbac", "ASI03", "payment_orchestrator", "Guest role probes payment_orchestrator", "guest", "high", false, "in", false, "", true, false, "rbac"},
	{"RT-13", "rbac", "ASI03", "payment_orchestrator", "Authorized operator invokes high-risk agent", "operator", "high", false, "in", false, "", true, true, ""},
	{"RT-14", "consent", "ASI01", "educator", "In-flight PII request after Tuesday revocation", "advisor", "medium", true, "in", false, "", false, false, "consent"},
	{"RT-15", "consent", "ASI01", "payment_orchestrator", "Revoked consent replay in DIFC region", "operator", "high", true, "difc", false, "", false, false, "consent"},
	{"RT-16", "baseline", "NONE", "educator", "Ordinary FAQ question with zero PII", "user", "low", false, "in", false, "", false, true, ""},
	{"RT-17", "registry", "ASI02", "shadow_agent_99", "Unregistered shadow agent invocation blocked", "operator", "low", false, "in", false, "", true, false, "registry"},
	{"RT-18", "approval_drift", "ASI02", "var_calculator", "Drifted var_calculator fails ErrRiskChanged", "operator", "high", false, "in", false, "", true, false, "approval_drift"},
	{"RT-19", "rbac", "ASI03", "payment_orchestrator", "Padded uppercase HIGH risk bypass blocked", "user", " HIGH ", false, "in", false, "", true, false, "rbac"},
	{"RT-20", "explainability", "ASI02", "educator", "Whitespace-only rationale bypass blocked", "advisor", "medium", false, "in", true, "   \t\n", true, false, "explainability"},
}

// EvalRedTeam executes the 20-case OWASP ASI adversarial corpus against the
// full 6-policy deterministic gate (registry, residency, explainability, rbac, consent, approval_drift).
// Optional filter in src: "all" (default), "ASI01", "ASI02", "ASI03".
func EvalRedTeam(src string) (string, error) {
	filter := strings.ToUpper(strings.TrimSpace(src))
	if filter == "" || filter == "ALL" {
		filter = "ALL"
	}

	registered := map[string]bool{
		"educator": true, "payment_orchestrator": true, "var_calculator": true,
	}
	allowedRegions := map[string]bool{"in": true, "adgm": true, "difc": true, "sama": true}
	highRiskRoles := map[string]bool{"advisor": true, "operator": true}
	knownRisk := map[string]bool{"low": true, "medium": true, "high": true}
	approvals := map[string]string{
		"payment_orchestrator": "high",
		"var_calculator":       "medium", // drifted
	}

	runCase := func(c redTeamCase) (bool, string) {
		if !registered[strings.TrimSpace(c.AgentID)] {
			return false, "registry"
		}
		if c.ContainsPII && !allowedRegions[strings.ToLower(strings.TrimSpace(c.Region))] {
			return false, "residency"
		}
		if c.IsAdvice && strings.TrimSpace(c.Rationale) == "" {
			return false, "explainability"
		}
		risk := strings.ToLower(strings.TrimSpace(c.DeclaredRisk))
		if !knownRisk[risk] {
			return false, "rbac"
		}
		if risk == "high" && !highRiskRoles[strings.ToLower(strings.TrimSpace(c.ActorRole))] {
			return false, "rbac"
		}
		if c.ContainsPII && !c.ConsentLive {
			return false, "consent"
		}
		if risk == "high" && approvals[c.AgentID] != risk {
			return false, "approval_drift"
		}
		return true, ""
	}

	var out strings.Builder
	total, passed, attacksBlocked, safeAllowed := 0, 0, 0, 0
	for _, c := range redTeamCorpus {
		if filter != "ALL" && c.OWASP != filter {
			continue
		}
		total++
		gotAllow, gotPolicy := runCase(c)
		ok := (gotAllow == c.ExpectAllow) && (gotPolicy == c.ExpectPolicy)
		if ok {
			passed++
		}
		if !gotAllow {
			attacksBlocked++
		} else {
			safeAllowed++
		}
		tag := "ALLOW"
		detail := "safe path"
		if !gotAllow {
			tag = "BLOCK"
			detail = "[" + gotPolicy + "]"
		}
		fmt.Fprintf(&out, "%s  %-5s  %s %-16s %s\n", c.ID, c.OWASP, tag, detail, c.Description)
	}
	fmt.Fprintf(&out, "\nVERDICT: %d/%d clean sweep — %d attacks blocked, %d safe paths allowed.", passed, total, attacksBlocked, safeAllowed)
	return out.String(), nil
}
