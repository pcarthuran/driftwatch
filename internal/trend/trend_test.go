package trend_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/trend"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
var t1 = time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
var t2 = time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)

func sampleEntries() []baseline.Entry {
	return []baseline.Entry{
		{
			SavedAt: t1,
			Results: []drift.Result{
				{ResourceID: "a", Status: "missing"},
				{ResourceID: "b", Status: "modified"},
			},
		},
		{
			SavedAt: t0,
			Results: []drift.Result{
				{ResourceID: "c", Status: "extra"},
			},
		},
		{
			SavedAt: t2,
			Results: []drift.Result{},
		},
	}
}

func TestCompute_Counts(t *testing.T) {
	r := trend.Compute(sampleEntries())
	if len(r.Points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(r.Points))
	}
	// Points should be sorted ascending by timestamp.
	if !r.Points[0].Timestamp.Equal(t0) {
		t.Errorf("expected first point at t0, got %v", r.Points[0].Timestamp)
	}
	if r.Points[0].Extra != 1 || r.Points[0].Total != 1 {
		t.Errorf("unexpected counts for t0 point: %+v", r.Points[0])
	}
	if r.Points[1].Missing != 1 || r.Points[1].Modified != 1 || r.Points[1].Total != 2 {
		t.Errorf("unexpected counts for t1 point: %+v", r.Points[1])
	}
	if r.Points[2].Total != 0 {
		t.Errorf("expected zero total for t2 point, got %d", r.Points[2].Total)
	}
}

func TestCompute_Empty(t *testing.T) {
	r := trend.Compute(nil)
	if len(r.Points) != 0 {
		t.Errorf("expected no points for empty input")
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	r := trend.Compute(sampleEntries())
	var buf bytes.Buffer
	if err := trend.Write(&buf, r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"Timestamp", "Missing", "Extra", "Modified", "Total"} {
		if !strings.Contains(out, col) {
			t.Errorf("output missing column %q", col)
		}
	}
}

func TestWrite_EmptyReport(t *testing.T) {
	var buf bytes.Buffer
	if err := trend.Write(&buf, trend.Report{}); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "No trend data") {
		t.Errorf("expected no-data message, got: %s", buf.String())
	}
}
