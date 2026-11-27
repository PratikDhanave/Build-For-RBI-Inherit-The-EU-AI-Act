// Repository verification CLI in pure Go — validates machine-readable cross-walks,
// sample JSON exports, policy YAML parity, Mermaid diagrams, and web-deck assets.
//
// Run from code/:
//
//	go run ./cmd/verify --root ..
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	root := flag.String("root", "..", "repository root directory")
	checkMermaid := flag.Bool("mermaid", true, "validate Mermaid blocks across README.md and docs/")
	checkAssets := flag.Bool("assets", true, "validate assets/* references in index.html")
	flag.Parse()

	must(validateCrosswalk(*root))
	must(validatePolicyParity(*root))
	if *checkMermaid {
		must(validateMermaid(*root))
	}
	if *checkAssets {
		must(validateDeckAssets(*root))
	}
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FAIL] %v\n", err)
		os.Exit(1)
	}
}

func validateCrosswalk(root string) error {
	raw, err := os.ReadFile(filepath.Join(root, "crosswalk/rbi-eu-ai-act.json"))
	if err != nil {
		return err
	}
	var doc struct {
		Recs []struct {
			Owner string `json:"owner"`
		} `json:"rbi_free_ai_recommendations"`
		Crosswalk []any `json:"rbi_to_eu_crosswalk"`
		MRM       []any `json:"rbi_draft_mrm_scorecard"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("parse rbi-eu-ai-act.json: %w", err)
	}
	if len(doc.Recs) != 26 {
		return fmt.Errorf("expected 26 RBI FREE-AI recommendations, got %d", len(doc.Recs))
	}
	reCount := 0
	for _, r := range doc.Recs {
		if r.Owner == "regulated_entity" {
			reCount++
		}
	}
	if reCount != 17 {
		return fmt.Errorf("expected 17 regulated_entity recommendations, got %d", reCount)
	}
	if len(doc.Crosswalk) != 7 {
		return fmt.Errorf("expected 7 crosswalk rows, got %d", len(doc.Crosswalk))
	}
	if len(doc.MRM) != 6 {
		return fmt.Errorf("expected 6 MRM scorecard rows, got %d", len(doc.MRM))
	}

	csvFile, err := os.Open(filepath.Join(root, "crosswalk/rbi-eu-ai-act.csv"))
	if err != nil {
		return err
	}
	defer csvFile.Close()
	rows, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return fmt.Errorf("parse rbi-eu-ai-act.csv: %w", err)
	}
	if len(rows)-1 != 7 {
		return fmt.Errorf("expected 7 CSV data rows, got %d", len(rows)-1)
	}

	for _, rel := range []string{
		"examples/annexure-vi-sample.json",
		"examples/inventory-export.json",
		"code/redteam/corpus.json",
	} {
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return err
		}
		var tmp any
		if err := json.Unmarshal(b, &tmp); err != nil {
			return fmt.Errorf("invalid JSON in %s: %w", rel, err)
		}
	}

	fmt.Println("[OK] Cross-walk, sample exports & red-team JSON verified (Go).")
	return nil
}

func validatePolicyParity(root string) error {
	a, err := os.ReadFile(filepath.Join(root, "code/config/policy.yaml"))
	if err != nil {
		return err
	}
	b, err := os.ReadFile(filepath.Join(root, "examples/policy.yaml"))
	if err != nil {
		return err
	}
	if !bytes.Equal(a, b) {
		return fmt.Errorf("code/config/policy.yaml and examples/policy.yaml differ")
	}
	fmt.Println("[OK] Board policy YAML parity verified (code/config == examples).")
	return nil
}

var (
	reMermaidBlock = regexp.MustCompile("(?s)```mermaid\\n(.*?)```")
	reInitJSON     = regexp.MustCompile(`^%%\{init:\s*(\{[^\n]+\})\s*\}%%`)
)

func validateMermaid(root string) error {
	files := []string{filepath.Join(root, "README.md")}
	err := filepath.WalkDir(filepath.Join(root, "docs"), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return err
	}

	total := 0
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		for _, m := range reMermaidBlock.FindAllStringSubmatch(string(raw), -1) {
			total++
			body := m[1]
			im := reInitJSON.FindStringSubmatch(body)
			if im == nil {
				return fmt.Errorf("expected single-line init JSON in %s", f)
			}
			var js any
			if err := json.Unmarshal([]byte(im[1]), &js); err != nil {
				return fmt.Errorf("invalid Mermaid init JSON in %s: %w", f, err)
			}
			if strings.Contains(body, "sequenceDiagram") {
				for _, line := range strings.Split(body, "\n") {
					trimmed := strings.TrimSpace(line)
					if !strings.HasPrefix(trimmed, "%%") && strings.Contains(line, ";") {
						return fmt.Errorf("semicolon breaks sequenceDiagram in %s: %s", f, line)
					}
				}
			}
		}
	}
	fmt.Printf("[OK] %d Mermaid diagrams validated across Markdown docs (Go).\n", total)
	return nil
}

var reAssetRef = regexp.MustCompile(`(?:src|href)="(assets/[^"]+)"`)

func validateDeckAssets(root string) error {
	raw, err := os.ReadFile(filepath.Join(root, "index.html"))
	if err != nil {
		return err
	}
	matches := reAssetRef.FindAllStringSubmatch(string(raw), -1)
	for _, m := range matches {
		ref := m[1]
		if _, err := os.Stat(filepath.Join(root, ref)); err != nil {
			return fmt.Errorf("missing deck asset %s: %w", ref, err)
		}
	}
	for _, wasmAsset := range []string{"assets/wasm_exec.js", "assets/labs.wasm"} {
		if _, err := os.Stat(filepath.Join(root, wasmAsset)); err != nil {
			return fmt.Errorf("missing Go WebAssembly asset %s: %w", wasmAsset, err)
		}
	}
	fmt.Printf("[OK] index.html asset references (%d) + Go WebAssembly assets verified.\n", len(matches))
	return nil
}
