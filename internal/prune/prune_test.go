package prune_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/prune"
)

func sampleResults() []drift.ResourceResult {
	return []drift.ResourceResult{
		{ID: "res-1", Provider: "aws", Type: "ec2", Drifted: false},
		{ID: "res-2", Provider: "aws", Type: "s3", Drifted: true},
		{ID: "res-3", Provider: "gcp", Type: "vm", Drifted: false},
	}
}

func makeHistory(ids []string, runs int) [][]drift.ResourceResult {
	var history [][]drift.ResourceResult
	for i := 0; i < runs; i++ {
		var run []drift.ResourceResult
		for _, id := range ids {
			run = append(run, drift.ResourceResult{ID: id, Drifted: false})
		}
		history = append(history, run)
	}
	return history
}

func TestPrune_NoHistory_RetainsAll(t *testing.T) {
	r := prune.Prune(sampleResults(), nil, 3)
	if len(r.Retained) != 3 {
		t.Fatalf("expected 3 retained, got %d", len(r.Retained))
	}
	if len(r.Pruned) != 0 {
		t.Fatalf("expected 0 pruned, got %d", len(r.Pruned))
	}
}

func TestPrune_DriftedAlwaysRetained(t *testing.T) {
	history := makeHistory([]string{"res-1", "res-2", "res-3"}, 5)
	r := prune.Prune(sampleResults(), history, 3)
	for _, res := range r.Retained {
		if res.ID == "res-2" {
			return
		}
	}
	t.Fatal("drifted resource res-2 should always be retained")
}

func TestPrune_CleanResourcePruned(t *testing.T) {
	history := makeHistory([]string{"res-1", "res-3"}, 4)
	r := prune.Prune(sampleResults(), history, 3)
	prunedIDs := make(map[string]bool)
	for _, res := range r.Pruned {
		prunedIDs[res.ID] = true
	}
	if !prunedIDs["res-1"] {
		t.Error("expected res-1 to be pruned")
	}
	if !prunedIDs["res-3"] {
		t.Error("expected res-3 to be pruned")
	}
	if prunedIDs["res-2"] {
		t.Error("res-2 is drifted and must not be pruned")
	}
}

func TestPrune_InsufficientHistory_Retains(t *testing.T) {
	history := makeHistory([]string{"res-1"}, 2)
	r := prune.Prune(sampleResults(), history, 5)
	for _, res := range r.Pruned {
		if res.ID == "res-1" {
			t.Fatal("res-1 should not be pruned with only 2 clean runs vs threshold 5")
		}
	}
}

func TestPrune_ZeroThreshold_RetainsAll(t *testing.T) {
	history := makeHistory([]string{"res-1", "res-3"}, 10)
	r := prune.Prune(sampleResults(), history, 0)
	if len(r.Pruned) != 0 {
		t.Fatalf("zero threshold should retain all, got %d pruned", len(r.Pruned))
	}
}

func TestWrite_ContainsSummary(t *testing.T) {
	result := prune.Result{
		Retained: sampleResults()[:2],
		Pruned:   sampleResults()[2:],
	}
	var buf bytes.Buffer
	if err := prune.Write(&buf, result); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Prune Summary") {
		t.Error("output missing 'Prune Summary' header")
	}
	if !strings.Contains(out, "res-3") {
		t.Error("output missing pruned resource res-3")
	}
}
