package tags_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/tags"
)

func sampleResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{ResourceID: "res-1", Status: drift.StatusModified},
		{ResourceID: "res-2", Status: drift.StatusMissing},
		{ResourceID: "res-3", Status: drift.StatusExtra},
	}
}

func sampleTagMap() map[string]map[string]string {
	return map[string]map[string]string{
		"res-1": {"env": "prod", "team": "platform"},
		"res-2": {"env": "staging", "team": "platform"},
		"res-3": {"env": "prod", "team": "data"},
	}
}

func TestNewFilter_KeyValue(t *testing.T) {
	f, err := tags.NewFilter([]string{"env=prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Required["env"] != "prod" {
		t.Errorf("expected Required[env]=prod, got %q", f.Required["env"])
	}
}

func TestNewFilter_KeyOnly(t *testing.T) {
	f, err := tags.NewFilter([]string{"team"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Keys) != 1 || f.Keys[0] != "team" {
		t.Errorf("expected Keys=[team], got %v", f.Keys)
	}
}

func TestNewFilter_EmptyKey_ReturnsError(t *testing.T) {
	_, err := tags.NewFilter([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestApply_NoFilter_ReturnsAll(t *testing.T) {
	results := sampleResults()
	out := tags.Apply(results, sampleTagMap(), nil)
	if len(out) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_FilterByKeyValue(t *testing.T) {
	f, _ := tags.NewFilter([]string{"env=prod"})
	out := tags.Apply(sampleResults(), sampleTagMap(), f)
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.ResourceID == "res-2" {
			t.Errorf("res-2 should have been filtered out")
		}
	}
}

func TestApply_FilterByMultipleTags(t *testing.T) {
	f, _ := tags.NewFilter([]string{"env=prod", "team=platform"})
	out := tags.Apply(sampleResults(), sampleTagMap(), f)
	if len(out) != 1 || out[0].ResourceID != "res-1" {
		t.Errorf("expected only res-1, got %v", out)
	}
}

func TestApply_FilterByKeyPresence(t *testing.T) {
	f, _ := tags.NewFilter([]string{"team"})
	out := tags.Apply(sampleResults(), sampleTagMap(), f)
	if len(out) != 3 {
		t.Errorf("expected 3 results (all have team tag), got %d", len(out))
	}
}

func TestApply_NoMatch_ReturnsEmpty(t *testing.T) {
	f, _ := tags.NewFilter([]string{"env=dev"})
	out := tags.Apply(sampleResults(), sampleTagMap(), f)
	if len(out) != 0 {
		t.Errorf("expected 0 results, got %d", len(out))
	}
}
