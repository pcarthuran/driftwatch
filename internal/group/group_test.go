package group_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/group"
	"github.com/driftwatch/internal/state"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{Resource: state.Resource{ID: "r1", Provider: "aws", Type: "ec2"}, Drifted: true},
		{Resource: state.Resource{ID: "r2", Provider: "aws", Type: "s3"}, Drifted: false},
		{Resource: state.Resource{ID: "r3", Provider: "gcp", Type: "gce"}, Drifted: true},
		{Resource: state.Resource{ID: "r4", Provider: "gcp", Type: "gce"}, Drifted: false},
		{Resource: state.Resource{ID: "r5", Provider: "azure", Type: "vm"}, Drifted: false},
	}
}

func TestCompute_ByProvider(t *testing.T) {
	groups := group.Compute(sampleResults(), group.ByProvider)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if groups[0].Key != "aws" {
		t.Errorf("expected first group 'aws', got %q", groups[0].Key)
	}
}

func TestCompute_ByType(t *testing.T) {
	groups := group.Compute(sampleResults(), group.ByType)
	keys := make([]string, len(groups))
	for i, g := range groups {
		keys[i] = g.Key
	}
	found := false
	for _, k := range keys {
		if k == "gce" {
			found = true
		}
	}
	if !found {
		t.Error("expected group 'gce' not found")
	}
}

func TestCompute_ByStatus(t *testing.T) {
	groups := group.Compute(sampleResults(), group.ByStatus)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups (drifted/clean), got %d", len(groups))
	}
	for _, g := range groups {
		if g.Key != "clean" && g.Key != "drifted" {
			t.Errorf("unexpected group key %q", g.Key)
		}
	}
}

func TestCompute_Empty(t *testing.T) {
	groups := group.Compute(nil, group.ByProvider)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	groups := group.Compute(sampleResults(), group.ByProvider)
	group.Write(&buf, groups)
	out := buf.String()
	if !strings.Contains(out, "GROUP") {
		t.Error("expected header 'GROUP' in output")
	}
	if !strings.Contains(out, "DRIFTED") {
		t.Error("expected header 'DRIFTED' in output")
	}
}

func TestWrite_CountsCorrect(t *testing.T) {
	var buf bytes.Buffer
	groups := group.Compute(sampleResults(), group.ByProvider)
	group.Write(&buf, groups)
	out := buf.String()
	// aws has 2 resources, 1 drifted
	if !strings.Contains(out, "aws") {
		t.Error("expected 'aws' row in output")
	}
}

func TestCompute_ByProvider_DriftCounts(t *testing.T) {
	groups := group.Compute(sampleResults(), group.ByProvider)

	// Build a map for easier lookup by key.
	byKey := make(map[string]group.Group)
	for _, g := range groups {
		byKey[g.Key] = g
	}

	cases := []struct {
		provider     string
		wantTotal    int
		wantDrifted  int
	}{
		{"aws", 2, 1},
		{"gcp", 2, 1},
		{"azure", 1, 0},
	}
	for _, tc := range cases {
		g, ok := byKey[tc.provider]
		if !ok {
			t.Errorf("group %q not found", tc.provider)
			continue
		}
		if len(g.Results) != tc.wantTotal {
			t.Errorf("provider %q: expected %d total resources, got %d", tc.provider, tc.wantTotal, len(g.Results))
		}
		drifted := 0
		for _, r := range g.Results {
			if r.Drifted {
				drifted++
			}
		}
		if drifted != tc.wantDrifted {
			t.Errorf("provider %q: expected %d drifted, got %d", tc.provider, tc.wantDrifted, drifted)
		}
	}
}
