package flatten_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/diff"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/flatten"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{
			Provider:     "aws",
			ResourceID:   "i-001",
			ResourceType: "ec2",
			Diffs: []diff.FieldDiff{
				{Field: "instance_type", Wanted: "t3.micro", Got: "t3.small", Kind: diff.KindModified},
			},
		},
		{
			Provider:     "aws",
			ResourceID:   "i-002",
			ResourceType: "ec2",
			Diffs:        nil,
		},
		{
			Provider:     "gcp",
			ResourceID:   "vm-abc",
			ResourceType: "compute",
			Diffs: []diff.FieldDiff{
				{Field: "zone", Wanted: "us-east1-b", Got: "us-west1-a", Kind: diff.KindModified},
				{Field: "machine_type", Wanted: "n1-standard-1", Got: "n1-standard-2", Kind: diff.KindModified},
			},
		},
	}
}

func TestFlatten_CountsRows(t *testing.T) {
	rows := flatten.Flatten(sampleResults())
	// 1 diff row + 1 ok row + 2 diff rows = 4
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}
}

func TestFlatten_OkRowHasNoField(t *testing.T) {
	rows := flatten.Flatten(sampleResults())
	var okRow *flatten.Row
	for i := range rows {
		if rows[i].ResourceID == "i-002" {
			okRow = &rows[i]
			break
		}
	}
	if okRow == nil {
		t.Fatal("expected ok row for i-002")
	}
	if okRow.Status != "ok" {
		t.Errorf("expected status ok, got %s", okRow.Status)
	}
	if okRow.Field != "" {
		t.Errorf("expected empty field for ok row, got %s", okRow.Field)
	}
}

func TestFlatten_SortedByProviderThenID(t *testing.T) {
	rows := flatten.Flatten(sampleResults())
	if rows[0].Provider != "aws" {
		t.Errorf("expected first provider aws, got %s", rows[0].Provider)
	}
	if rows[len(rows)-1].Provider != "gcp" {
		t.Errorf("expected last provider gcp, got %s", rows[len(rows)-1].Provider)
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	rows := flatten.Flatten(sampleResults())
	var buf bytes.Buffer
	if err := flatten.Write(&buf, rows); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "PROVIDER") {
		t.Error("expected header PROVIDER in output")
	}
	if !strings.Contains(output, "WANTED") {
		t.Error("expected header WANTED in output")
	}
}

func TestWrite_ContainsResourceIDs(t *testing.T) {
	rows := flatten.Flatten(sampleResults())
	var buf bytes.Buffer
	_ = flatten.Write(&buf, rows)
	output := buf.String()
	for _, id := range []string{"i-001", "i-002", "vm-abc"} {
		if !strings.Contains(output, id) {
			t.Errorf("expected resource ID %s in output", id)
		}
	}
}
