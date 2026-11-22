package registry

import (
	"testing"
	"time"
)

func TestLiveRegistryZeroDriftExport(t *testing.T) {
	reg := NewLiveRegistry()
	reg.Register(AgentDescriptor{
		ID:             "payment_orchestrator",
		Domain:         "payments",
		DeclaredRisk:   "high",
		Owner:          "payments-platform@bank",
		VendorModel:    "deterministic+llm-v2",
		DataCategories: []string{"transaction", "account_id"},
		DisclosureMode: "first-event",
	})
	reg.Register(AgentDescriptor{
		ID:             "educator",
		Domain:         "financial-literacy",
		DeclaredRisk:   "low",
		Owner:          "retail-experience@bank",
		VendorModel:    "local-slm",
		DataCategories: []string{"public-faq"},
		DisclosureMode: "first-event",
	})

	exp := reg.Export(time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC))
	if exp.TotalLiveAgents != 2 {
		t.Fatalf("expected 2 live agents, got %d", exp.TotalLiveAgents)
	}
	if exp.Agents[0].ID != "educator" || exp.Agents[1].ID != "payment_orchestrator" {
		t.Fatalf("expected deterministic sorted order, got %+v", exp.Agents)
	}
	if _, err := exp.ToJSON(); err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
}
