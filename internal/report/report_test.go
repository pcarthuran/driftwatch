package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/report"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "server-1", Status: drift.StatusMissing, Detail: "not found in live state"},
		{ResourceID: "server-2", Status: drift.StatusModified, Detail: "cpu: want 2, got 4"},
	}
}

func TestWrite_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.New(&buf, report.FormatText)
	if err := w.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWrite_TextWithDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.New(&buf, report.FormatText)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "server-1") {
		t.Errorf("expected server-1 in output, got: %s", out)
	}
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING status in output, got: %s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := report.New(&buf, report.FormatJSON)
	results := sampleResults()
	if err := w.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []drift.Result
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(decoded))
	}
}

func TestWrite_JSONNoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.New(&buf, report.FormatJSON)
	if err := w.Write([]drift.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []drift.Result
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded) != 0 {
		t.Errorf("expected empty results, got %d", len(decoded))
	}
}
