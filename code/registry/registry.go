package registry

import (
	"encoding/json"
	"sort"
	"time"
)

// AgentDescriptor captures the live metadata of a registered agent.
// Generating the inventory directly from this runtime registry guarantees zero drift
// (RBI FREE-AI Rec 23 / EU AI Act Annex IV).
type AgentDescriptor struct {
	ID             string   `json:"id"`
	Domain         string   `json:"domain"`
	DeclaredRisk   string   `json:"declared_risk"`
	Owner          string   `json:"owner"`
	VendorModel    string   `json:"vendor_model"`
	DataCategories []string `json:"data_categories"`
	DisclosureMode string   `json:"disclosure_mode"`
}

// InventoryExport is the structured compliance pack generated from the live registry.
type InventoryExport struct {
	GeneratedAt        string            `json:"generated_at"`
	Source             string            `json:"source"`
	InstrumentMappings []string          `json:"instrument_mappings"`
	TotalLiveAgents    int               `json:"total_live_agents"`
	Agents             []AgentDescriptor `json:"agents"`
}

type LiveRegistry struct {
	agents map[string]AgentDescriptor
}

func NewLiveRegistry() *LiveRegistry {
	return &LiveRegistry{agents: make(map[string]AgentDescriptor)}
}

func (r *LiveRegistry) Register(a AgentDescriptor) {
	r.agents[a.ID] = a
}

func (r *LiveRegistry) Export(now time.Time) InventoryExport {
	ids := make([]string, 0, len(r.agents))
	for id := range r.agents {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	list := make([]AgentDescriptor, 0, len(ids))
	for _, id := range ids {
		list = append(list, r.agents[id])
	}
	return InventoryExport{
		GeneratedAt:        now.UTC().Format(time.RFC3339),
		Source:             "live-runtime-registry (zero-drift)",
		InstrumentMappings: []string{"RBI FREE-AI Rec 23", "EU AI Act Annex IV", "RBI Draft MRM Full Inventory"},
		TotalLiveAgents:    len(list),
		Agents:             list,
	}
}

func (e InventoryExport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}
