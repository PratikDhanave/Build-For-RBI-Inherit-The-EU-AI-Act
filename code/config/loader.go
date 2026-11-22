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
	HighRiskRoles              map[string]bool
	RequireRationaleFor        []string
	EmitDisclosureOnFirstEvent bool
	AuditHashAlgorithm         string
}

// LoadPolicyYAML parses our board-approved policy.yaml using zero external dependencies
// so the entire governance suite stays 100% offline-runnable.
func LoadPolicyYAML(path string) (PolicyConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return PolicyConfig{}, fmt.Errorf("read policy yaml: %w", err)
	}
	cfg := PolicyConfig{
		PIIAllowedRegions: make(map[string]bool),
		HighRiskRoles:     make(map[string]bool),
	}

	section := ""
	listKey := ""
	for _, line := range strings.Split(string(raw), "\n") {
		// Strip inline comments outside quotes
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if strings.HasPrefix(trimmed, "- ") {
			val := unquote(strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
			switch listKey {
			case "pii_allowed_regions":
				cfg.PIIAllowedRegions[val] = true
			case "high_risk_allowed_roles":
				cfg.HighRiskRoles[val] = true
			case "require_rationale_for":
				cfg.RequireRationaleFor = append(cfg.RequireRationaleFor, val)
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
			cfg.EmitDisclosureOnFirstEvent = (val == "true")
		case section == "audit" && key == "hash_algorithm":
			cfg.AuditHashAlgorithm = val
		}
	}

	if cfg.ApprovedBy == "" || len(cfg.PIIAllowedRegions) == 0 {
		return PolicyConfig{}, fmt.Errorf("policy yaml %s missing required governance fields", path)
	}
	return cfg, nil
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}
