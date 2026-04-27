package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/filter"
	"github.com/driftwatch/internal/state"
)

func sampleResources() []state.Resource {
	return []state.Resource{
		{ID: "r1", Type: "instance", Provider: "aws", Attributes: map[string]interface{}{"env": "prod"}},
		{ID: "r2", Type: "bucket", Provider: "aws", Attributes: map[string]interface{}{"env": "dev"}},
		{ID: "r3", Type: "instance", Provider: "gcp", Attributes: map[string]interface{}{"env": "prod"}},
		{ID: "r4", Type: "disk", Provider: "azure", Attributes: map[string]interface{}{"env": "staging"}},
	}
}

func TestApply_NoFilter(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{})
	if len(res) != 4 {
		t.Fatalf("expected 4 resources, got %d", len(res))
	}
}

func TestApply_FilterByProvider(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{Providers: []string{"aws"}})
	if len(res) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(res))
	}
	for _, r := range res {
		if r.Provider != "aws" {
			t.Errorf("unexpected provider %s", r.Provider)
		}
	}
}

func TestApply_FilterByType(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{Types: []string{"instance"}})
	if len(res) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(res))
	}
}

func TestApply_FilterByID(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{IDs: []string{"r1", "r4"}})
	if len(res) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(res))
	}
}

func TestApply_FilterByLabel(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{LabelKey: "env", LabelVal: "prod"})
	if len(res) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(res))
	}
}

func TestApply_FilterByLabelKeyOnly(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{LabelKey: "env"})
	if len(res) != 4 {
		t.Fatalf("expected 4 resources, got %d", len(res))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{
		Providers: []string{"aws"},
		Types:     []string{"instance"},
	})
	if len(res) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(res))
	}
	if res[0].ID != "r1" {
		t.Errorf("expected r1, got %s", res[0].ID)
	}
}

func TestApply_NoMatch(t *testing.T) {
	res := filter.Apply(sampleResources(), filter.Options{Providers: []string{"unknown"}})
	if len(res) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(res))
	}
}
