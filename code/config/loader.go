package config

import (
	"fmt"
	"os"
	"strings"
)

// PolicyConfig holds the board-approved declarative governance rules loaded
// from policy.yaml (RBI FREE-AI Rec 14 / EU AI Act Article 17).
type PolicyConfig struct {
	Version                    string
	ApprovedBy                 string
	EffectiveDate              string
	PIIAllowedRegions          map[string]bool
	DenyRegions                map[string]bool
	HighRiskRoles              map[string]bool
	RequireRationaleFor        []string
	EmitDisclosureOnFirstEvent bool
	AuditHashAlgorithm         string
}

// LoadPolicyYAML parses our board-approved policy.yaml using zero external dependencies
// so the entire governance suite stays 100% offline-runnable while strictly validating
// governance invariants (quote-aware comments, explicit deny precedence, no tab indentation).
func LoadPolicyYAML(path string) (PolicyConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return PolicyConfig{}, fmt.Errorf("read policy yaml: %w", err)
	}
	cfg := PolicyConfig{
		PIIAllowedRegions: make(map[string]bool),
		DenyRegions:       make(map[string]bool),
		HighRiskRoles:     make(map[string]bool),
	}

	section := ""
	listKey := ""
	for lineNum, origLine := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(origLine, "\t") || strings.Contains(leadingWhitespace(origLine), "\t") {
			return PolicyConfig{}, fmt.Errorf("policy yaml %s:%d uses tab indentation (use spaces)", path, lineNum+1)
		}
		line := stripCommentOutsideQuotes(origLine)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if strings.HasPrefix(trimmed, "- ") {
			val := strings.ToLower(unquote(strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))))
			switch listKey {
			case "pii_allowed_regions":
				if val != "" {
					cfg.PIIAllowedRegions[val] = true
				}
			case "deny_regions":
				if val != "" {
					cfg.DenyRegions[val] = true
				}
			case "high_risk_allowed_roles":
				if val != "" {
					cfg.HighRiskRoles[val] = true
				}
			case "require_rationale_for":
				if val != "" {
					cfg.RequireRationaleFor = append(cfg.RequireRationaleFor, val)
				}
			}
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := unquote(strings.TrimSpace(parts[1]))

		if indent == 0 {
			section = key
			listKey = ""
			switch key {
			case "version":
				cfg.Version = val
			case "approved_by":
				cfg.ApprovedBy = val
			case "effective_date":
				cfg.EffectiveDate = val
			}
			continue
		}

		if val == "" {
			if indent == 2 {
				section = key
				listKey = ""
			} else {
				listKey = key
			}
			continue
		}

		listKey = ""
		switch {
		case section == "disclosure" && key == "emit_on_first_event":
			cfg.EmitDisclosureOnFirstEvent = (strings.ToLower(val) == "true")
		case section == "audit" && key == "hash_algorithm":
			cfg.AuditHashAlgorithm = strings.ToLower(val)
		}
	}

	// Explicit deny_regions always takes precedence over pii_allowed_regions.
	for denied := range cfg.DenyRegions {
		delete(cfg.PIIAllowedRegions, denied)
	}

	if cfg.ApprovedBy == "" || len(cfg.PIIAllowedRegions) == 0 {
		return PolicyConfig{}, fmt.Errorf("policy yaml %s missing required governance fields (approved_by / pii_allowed_regions)", path)
	}
	if len(cfg.HighRiskRoles) == 0 {
		return PolicyConfig{}, fmt.Errorf("policy yaml %s missing high_risk_allowed_roles", path)
	}
	if cfg.AuditHashAlgorithm != "sha256" {
		return PolicyConfig{}, fmt.Errorf("policy yaml %s unsupported audit hash_algorithm %q (must be sha256)", path, cfg.AuditHashAlgorithm)
	}
	return cfg, nil
}

func leadingWhitespace(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' {
			return s[:i]
		}
	}
	return s
}

func stripCommentOutsideQuotes(line string) string {
	inSingle := false
	inDouble := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if c == '\'' && !inDouble {
			inSingle = !inSingle
		} else if c == '"' && !inSingle {
			inDouble = !inDouble
		} else if c == '#' && !inSingle && !inDouble {
			return line[:i]
		}
	}
	return line
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}
