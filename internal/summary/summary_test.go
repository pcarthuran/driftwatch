package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/summary"
)

func sampleResults() []drift.DetectResult {
	return []drift.DetectResult{
		{ResourceID: "res-1", Provider: "aws", ResourceType: "ec2", Status: drift.StatusOK},
		{ResourceID: "res-2", Provider: "aws", ResourceType: "ec2", Status: drift.StatusModified},
		{ResourceID: "res-3", Provider: "gcp", ResourceType: "vm", Status: drift.StatusMissing},
		{ResourceID: "res-4", Provider: "azure", ResourceType: "vm", Status: drift.StatusExtra},
	}
}

func TestCompute_Counts(t *testing.T) {
	s := summary.Compute(sampleResults())

	if s.Total != 4 {
		t.Errorf("Total: got %d, want 4", s.Total)
	}
	if s.Clean != 1 {
		t.Errorf("Clean: got %d, want 1", s.Clean)
	}
	if s.Modified != 1 {
		t.Errorf("Modified: got %d, want 1", s.Modified)
	}
	if s.Missing != 1 {
		t.Errorf("Missing: got %d, want 1", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("Extra: got %d, want 1", s.Extra)
	}
}

func TestCompute_AllClean(t *testing.T) {
	results := []drift.DetectResult{
		{ResourceID: "a", Status: drift.StatusOK},
		{ResourceID: "b", Status: drift.StatusOK},
	}
	s := summary.Compute(results)
	if s.Clean != 2 || s.Modified != 0 || s.Missing != 0 || s.Extra != 0 {
		t.Errorf("unexpected stats for all-clean results: %+v", s)
	}
}

func TestWrite_ContainsSummaryHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := summary.Write(&buf, sampleResults()); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SUMMARY") {
		t.Error("output missing SUMMARY header")
	}
	if !strings.Contains(out, "Total resources:") {
		t.Error("output missing 'Total resources:' line")
	}
}

func TestWrite_DriftedSectionPresent(t *testing.T) {
	var buf bytes.Buffer
	if err := summary.Write(&buf, sampleResults()); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DRIFTED RESOURCES") {
		t.Error("expected DRIFTED RESOURCES section")
	}
	if !strings.Contains(out, "res-2") {
		t.Error("expected drifted resource res-2 in output")
	}
}

func TestWrite_NoDriftSection_WhenAllClean(t *testing.T) {
	results := []drift.DetectResult{
		{ResourceID: "x", Provider: "aws", ResourceType: "s3", Status: drift.StatusOK},
	}
	var buf bytes.Buffer
	if err := summary.Write(&buf, results); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "DRIFTED RESOURCES") {
		t.Error("did not expect DRIFTED RESOURCES section when all resources are clean")
	}
}
