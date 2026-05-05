package tags_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/tags"
)

// TestIntegration_ChainedFilters verifies that multiple tag expressions
// applied together narrow results correctly across a realistic resource set.
func TestIntegration_ChainedFilters(t *testing.T) {
	results := []drift.ResourceDiff{
		{ResourceID: "vm-prod-1", Status: drift.StatusModified},
		{ResourceID: "vm-prod-2", Status: drift.StatusModified},
		{ResourceID: "vm-staging-1", Status: drift.StatusMissing},
		{ResourceID: "db-prod-1", Status: drift.StatusExtra},
	}
	tagMap := map[string]map[string]string{
		"vm-prod-1":    {"env": "prod", "type": "vm", "region": "us-east-1"},
		"vm-prod-2":    {"env": "prod", "type": "vm", "region": "eu-west-1"},
		"vm-staging-1": {"env": "staging", "type": "vm", "region": "us-east-1"},
		"db-prod-1":    {"env": "prod", "type": "db", "region": "us-east-1"},
	}

	// Filter: env=prod AND type=vm
	f, err := tags.NewFilter([]string{"env=prod", "type=vm"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tags.Apply(results, tagMap, f)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.ResourceID != "vm-prod-1" && r.ResourceID != "vm-prod-2" {
			t.Errorf("unexpected resource %q in results", r.ResourceID)
		}
	}
}

// TestIntegration_KeyPresenceAcrossProviders ensures key-only filters
// work when the tag map is populated from mixed provider resources.
func TestIntegration_KeyPresenceAcrossProviders(t *testing.T) {
	results := []drift.ResourceDiff{
		{ResourceID: "a", Status: drift.StatusModified},
		{ResourceID: "b", Status: drift.StatusModified},
		{ResourceID: "c", Status: drift.StatusModified},
	}
	tagMap := map[string]map[string]string{
		"a": {"owner": "alice"},
		"b": {},
		"c": {"owner": "bob", "cost-center": "eng"},
	}

	f, err := tags.NewFilter([]string{"owner"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tags.Apply(results, tagMap, f)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.ResourceID == "b" {
			t.Errorf("resource b (no owner tag) should have been excluded")
		}
	}
}

// TestIntegration_EmptyExprs_PassesThrough confirms that an empty
// filter expression list returns the full result set unchanged.
func TestIntegration_EmptyExprs_PassesThrough(t *testing.T) {
	results := sampleResults()
	f, err := tags.NewFilter([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tags.Apply(results, sampleTagMap(), f)
	if len(out) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(out))
	}
}
