package classify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/classify"
	"github.com/driftwatch/internal/drift"
)

func sampleResults() []drift.ResourceResult {
	return []drift.ResourceResult{
		{ResourceID: "res-1", ResourceType: "instance", Provider: "aws", Status: drift.StatusMissing},
		{ResourceID: "res-2", ResourceType: "bucket", Provider: "gcp", Status: drift.StatusExtra},
		{
			ResourceID: "res-3", ResourceType: "disk", Provider: "azure",
			Status: drift.StatusModified,
			Diffs: []drift.FieldDiff{{Field: "a"}, {Field: "b"}, {Field: "c"}},
		},
		{
			ResourceID: "res-4", ResourceType: "vm", Provider: "aws",
			Status: drift.StatusModified,
			Diffs: []drift.FieldDiff{{Field: "x"}},
		},
		{ResourceID: "res-5", ResourceType: "subnet", Provider: "gcp", Status: drift.StatusOK},
	}
}

func TestClassify_Severities(t *testing.T) {
	results := sampleResults()
	got := classify.Classify(results)

	if len(got) != len(results) {
		t.Fatalf("expected %d classifications, got %d", len(results), len(got))
	}

	expected := []classify.Severity{
		classify.SeverityCritical,
		classify.SeverityHigh,
		classify.SeverityHigh,
		classify.SeverityLow,
		classify.SeverityNone,
	}
	for i, c := range got {
		if c.Severity != expected[i] {
			t.Errorf("result[%d]: expected severity %s, got %s", i, expected[i], c.Severity)
		}
	}
}

func TestClassify_Empty(t *testing.T) {
	got := classify.Classify(nil)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	classifications := classify.Classify(sampleResults())
	if err := classify.Write(&buf, classifications); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SEVERITY") {
		t.Error("expected output to contain SEVERITY header")
	}
	if !strings.Contains(out, "CRITICAL") {
		t.Error("expected output to contain CRITICAL row")
	}
}

func TestWrite_EmptyClassifications(t *testing.T) {
	var buf bytes.Buffer
	if err := classify.Write(&buf, nil); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SEVERITY") {
		t.Error("expected header even with no results")
	}
}
