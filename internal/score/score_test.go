package score_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/score"
)

func sampleResults(statuses ...drift.Status) []drift.ResourceResult {
	results := make([]drift.ResourceResult, len(statuses))
	for i, s := range statuses {
		results[i] = drift.ResourceResult{
			ResourceID: fmt.Sprintf("res-%d", i),
			Status:     s,
		}
	}
	return results
}

func TestCompute_Empty(t *testing.T) {
	r := score.Compute(nil)
	if r.Score != 100.0 {
		t.Errorf("expected 100.0, got %.2f", r.Score)
	}
	if r.Grade != "A" {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
}

func TestCompute_AllClean(t *testing.T) {
	results := []drift.ResourceResult{
		{ResourceID: "r1", Status: drift.StatusMatch},
		{ResourceID: "r2", Status: drift.StatusMatch},
	}
	r := score.Compute(results)
	if r.Score != 100.0 {
		t.Errorf("expected 100.0, got %.2f", r.Score)
	}
	if r.Clean != 2 || r.Drifted != 0 {
		t.Errorf("unexpected counts: clean=%d drifted=%d", r.Clean, r.Drifted)
	}
}

func TestCompute_MixedDrift(t *testing.T) {
	results := []drift.ResourceResult{
		{ResourceID: "r1", Status: drift.StatusMatch},
		{ResourceID: "r2", Status: drift.StatusMatch},
		{ResourceID: "r3", Status: drift.StatusMissing},
		{ResourceID: "r4", Status: drift.StatusModified},
	}
	r := score.Compute(results)
	if r.Total != 4 {
		t.Errorf("expected total 4, got %d", r.Total)
	}
	if r.Clean != 2 {
		t.Errorf("expected 2 clean, got %d", r.Clean)
	}
	if r.Drifted != 2 {
		t.Errorf("expected 2 drifted, got %d", r.Drifted)
	}
	if r.Score != 50.0 {
		t.Errorf("expected score 50.0, got %.2f", r.Score)
	}
	if r.Grade != "D" {
		t.Errorf("expected grade D, got %s", r.Grade)
	}
}

func TestCompute_GradeF(t *testing.T) {
	var results []drift.ResourceResult
	for i := 0; i < 10; i++ {
		results = append(results, drift.ResourceResult{ResourceID: fmt.Sprintf("r%d", i), Status: drift.StatusExtra})
	}
	r := score.Compute(results)
	if r.Grade != "F" {
		t.Errorf("expected grade F, got %s", r.Grade)
	}
}

func TestWrite_ContainsGrade(t *testing.T) {
	results := []drift.ResourceResult{
		{ResourceID: "r1", Status: drift.StatusMatch},
		{ResourceID: "r2", Status: drift.StatusMissing},
	}
	r := score.Compute(results)
	var buf bytes.Buffer
	if err := score.Write(&buf, r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Grade:") {
		t.Error("output missing Grade field")
	}
	if !strings.Contains(out, "Score:") {
		t.Error("output missing Score field")
	}
	if !strings.Contains(out, "Drifted:") {
		t.Error("output missing Drifted field")
	}
}
