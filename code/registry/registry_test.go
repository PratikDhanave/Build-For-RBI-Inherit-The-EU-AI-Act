package registry

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLiveRegistryZeroDriftExportAndRaceSafety(t *testing.T) {
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

	if !reg.Has("educator") || reg.Has("shadow_agent") {
		t.Fatalf("unexpected Has results: educator=%v shadow=%v", reg.Has("educator"), reg.Has("shadow_agent"))
	}

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

	// Concurrent Register + Has + Export must pass -race cleanly
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				reg.Register(AgentDescriptor{ID: fmt.Sprintf("agent-%d-%d", id, j), DeclaredRisk: "low"})
				_ = reg.Has("educator")
				_ = reg.Export(time.Now().UTC())
			}
		}(i)
	}
	wg.Wait()
}
