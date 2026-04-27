package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/export"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{
			ResourceID: "res-1",
			Provider:   "aws",
			Type:       "instance",
			Status:     drift.StatusModified,
			Diffs: []drift.FieldDiff{
				{Field: "size", Declared: "t2.micro", Live: "t2.large"},
			},
		},
		{
			ResourceID: "res-2",
			Provider:   "gcp",
			Type:       "bucket",
			Status:     drift.StatusMissing,
			Diffs:      nil,
		},
	}
}

func TestWrite_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(sampleResults(), export.Options{Format: export.FormatCSV, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "resource_id") {
		t.Error("expected CSV header row")
	}
	if !strings.Contains(output, "res-1") {
		t.Error("expected res-1 in output")
	}
	if !strings.Contains(output, "t2.micro") {
		t.Error("expected declared value in output")
	}
	if !strings.Contains(output, "t2.large") {
		t.Error("expected live value in output")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(sampleResults(), export.Options{Format: export.FormatJSON, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []drift.Result
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded) != 2 {
		t.Errorf("expected 2 results, got %d", len(decoded))
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(nil, export.Options{Format: "xml", Writer: &buf})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_CSVNoDrift(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "res-3", Provider: "azure", Type: "vm", Status: drift.StatusOK},
	}
	var buf bytes.Buffer
	err := export.Write(results, export.Options{Format: export.FormatCSV, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "res-3") {
		t.Error("expected res-3 in CSV output")
	}
}
