package rollup_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func sampleProviders() []rollup.ProviderResult {
	return []rollup.ProviderResult{
		{
			Provider: "aws",
			Results: []drift.Result{
				{ResourceID: "r1", Status: drift.StatusMatch},
				{ResourceID: "r2", Status: drift.StatusModified},
				{ResourceID: "r3", Status: drift.StatusMissing},
			},
		},
		{
			Provider: "gcp",
			Results: []drift.Result{
				{ResourceID: "g1", Status: drift.StatusMatch},
				{ResourceID: "g2", Status: drift.StatusExtra},
			},
		},
	}
}

func TestCompute_Counts(t *testing.T) {
	r := rollup.Compute(sampleProviders())

	if r.Total != 5 {
		t.Errorf("Total: want 5, got %d", r.Total)
	}
	if r.Clean != 2 {
		t.Errorf("Clean: want 2, got %d", r.Clean)
	}
	if r.Drifted != 3 {
		t.Errorf("Drifted: want 3, got %d", r.Drifted)
	}
	if r.Missing != 1 {
		t.Errorf("Missing: want 1, got %d", r.Missing)
	}
	if r.Extra != 1 {
		t.Errorf("Extra: want 1, got %d", r.Extra)
	}
	if r.Modified != 1 {
		t.Errorf("Modified: want 1, got %d", r.Modified)
	}
}

func TestCompute_Empty(t *testing.T) {
	r := rollup.Compute(nil)
	if r.Total != 0 || r.Drifted != 0 {
		t.Errorf("expected zero counts for empty input, got %+v", r)
	}
}

func TestWrite_ContainsProviders(t *testing.T) {
	var buf bytes.Buffer
	r := rollup.Compute(sampleProviders())
	if err := rollup.Write(&buf, r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"aws", "gcp", "PROVIDER", "TOTAL", "CLEAN", "DRIFTED", "Summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestWrite_SummaryLine(t *testing.T) {
	var buf bytes.Buffer
	r := rollup.Compute(sampleProviders())
	rollup.Write(&buf, r) //nolint:errcheck
	if !strings.Contains(buf.String(), "total=5") {
		t.Errorf("expected total=5 in summary line")
	}
}

func TestWrite_AllClean(t *testing.T) {
	providers := []rollup.ProviderResult{
		{
			Provider: "azure",
			Results: []drift.Result{
				{ResourceID: "a1", Status: drift.StatusMatch},
			},
		},
	}
	var buf bytes.Buffer
	r := rollup.Compute(providers)
	rollup.Write(&buf, r) //nolint:errcheck
	if !strings.Contains(buf.String(), "drifted=0") {
		t.Errorf("expected drifted=0 for all-clean input")
	}
}
