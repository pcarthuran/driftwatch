package export_test

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/export"
)

func TestCSV_RowCountMatchesDiffs(t *testing.T) {
	results := []drift.Result{
		{
			ResourceID: "r1",
			Provider:   "aws",
			Type:       "sg",
			Status:     drift.StatusModified,
			Diffs: []drift.FieldDiff{
				{Field: "port", Declared: "80", Live: "443"},
				{Field: "proto", Declared: "tcp", Live: "udp"},
			},
		},
	}

	var buf bytes.Buffer
	if err := export.Write(results, export.Options{Format: export.FormatCSV, Writer: &buf}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(strings.NewReader(buf.String()))
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	// 1 header + 2 diff rows
	if len(rows) != 3 {
		t.Errorf("expected 3 rows (header + 2 diffs), got %d", len(rows))
	}
}

func TestCSV_FieldsOrdered(t *testing.T) {
	results := []drift.Result{
		{
			ResourceID: "r2",
			Provider:   "gcp",
			Type:       "disk",
			Status:     drift.StatusModified,
			Diffs: []drift.FieldDiff{
				{Field: "zone", Declared: "us-east1", Live: "us-west1"},
				{Field: "size", Declared: "100", Live: "200"},
			},
		},
	}

	var buf bytes.Buffer
	if err := export.Write(results, export.Options{Format: export.FormatCSV, Writer: &buf}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(strings.NewReader(buf.String()))
	rows, _ := r.ReadAll()
	// rows[1] should be the alphabetically first field: "size"
	if len(rows) < 3 {
		t.Fatalf("expected at least 3 rows, got %d", len(rows))
	}
	if rows[1][4] != "size" {
		t.Errorf("expected first diff field to be 'size', got %q", rows[1][4])
	}
}

func TestJSON_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write([]drift.Result{}, export.Options{Format: export.FormatJSON, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if output != "[]" {
		t.Errorf("expected '[]' for empty results, got %q", output)
	}
}
