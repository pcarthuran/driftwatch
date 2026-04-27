package diff_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/diff"
)

func TestFieldDiff_String(t *testing.T) {
	f := diff.FieldDiff{Field: "region", Expected: "us-east-1", Actual: "eu-west-1"}
	got := f.String()
	if !strings.Contains(got, "region") {
		t.Errorf("expected field name in output, got: %s", got)
	}
	if !strings.Contains(got, "us-east-1") {
		t.Errorf("expected expected value in output, got: %s", got)
	}
	if !strings.Contains(got, "eu-west-1") {
		t.Errorf("expected actual value in output, got: %s", got)
	}
}

func TestResourceDiff_Summary_Missing(t *testing.T) {
	r := diff.ResourceDiff{ResourceID: "res-1", Kind: "missing"}
	got := r.Summary()
	if !strings.Contains(got, "[MISSING]") || !strings.Contains(got, "res-1") {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestResourceDiff_Summary_Extra(t *testing.T) {
	r := diff.ResourceDiff{ResourceID: "res-2", Kind: "extra"}
	got := r.Summary()
	if !strings.Contains(got, "[EXTRA]") || !strings.Contains(got, "res-2") {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestResourceDiff_Summary_Modified(t *testing.T) {
	r := diff.ResourceDiff{
		ResourceID: "res-3",
		Kind:       "modified",
		Fields: []diff.FieldDiff{
			{Field: "size", Expected: "small", Actual: "large"},
		},
	}
	got := r.Summary()
	if !strings.Contains(got, "[MODIFIED]") {
		t.Errorf("expected [MODIFIED] in summary, got: %s", got)
	}
	if !strings.Contains(got, "size") {
		t.Errorf("expected field name in summary, got: %s", got)
	}
}

func TestCompareFields_NoDiff(t *testing.T) {
	expected := map[string]interface{}{"region": "us-east-1", "type": "t2.micro"}
	actual := map[string]interface{}{"region": "us-east-1", "type": "t2.micro"}
	diffs := diff.CompareFields(expected, actual)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompareFields_ValueChanged(t *testing.T) {
	expected := map[string]interface{}{"region": "us-east-1"}
	actual := map[string]interface{}{"region": "eu-west-1"}
	diffs := diff.CompareFields(expected, actual)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "region" {
		t.Errorf("expected field 'region', got %q", diffs[0].Field)
	}
}

func TestCompareFields_MissingInActual(t *testing.T) {
	expected := map[string]interface{}{"region": "us-east-1", "size": "large"}
	actual := map[string]interface{}{"region": "us-east-1"}
	diffs := diff.CompareFields(expected, actual)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "size" {
		t.Errorf("expected field 'size', got %q", diffs[0].Field)
	}
}

func TestCompareFields_ExtraInActual(t *testing.T) {
	expected := map[string]interface{}{"region": "us-east-1"}
	actual := map[string]interface{}{"region": "us-east-1", "extra": "value"}
	diffs := diff.CompareFields(expected, actual)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "extra" {
		t.Errorf("expected field 'extra', got %q", diffs[0].Field)
	}
	if diffs[0].Expected != nil {
		t.Errorf("expected Expected to be nil for extra field")
	}
}
