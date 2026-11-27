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
