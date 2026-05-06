package lint_test

import (
	"testing"

	"github.com/driftwatch/internal/lint"
	"github.com/driftwatch/internal/state"
)

func TestIntegration_MultipleResourceErrors(t *testing.T) {
	snap := &state.Snapshot{
		Resources: []state.Resource{
			{ID: "", Type: "", Provider: "aws", Fields: map[string]any{"k": "v"}},
			{ID: "res-2", Type: "bucket", Provider: "", Fields: nil},
		},
	}
	result := lint.Run(snap)
	if !result.HasErrors() {
		t.Fatal("expected errors")
	}
	// first resource: missing id (error) + missing type (error) = 2
	// second resource: missing provider (warn) + empty fields (warn) = 2
	if len(result.Findings) != 4 {
		t.Fatalf("expected 4 findings, got %d", len(result.Findings))
	}
}

func TestIntegration_CleanSnapshot_NoFindings(t *testing.T) {
	snap := &state.Snapshot{
		Resources: []state.Resource{
			{ID: "a", Type: "vm", Provider: "gcp", Fields: map[string]any{"zone": "us-central1"}},
			{ID: "b", Type: "bucket", Provider: "aws", Fields: map[string]any{"region": "eu-west-1"}},
		},
	}
	result := lint.Run(snap)
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d: %v", len(result.Findings), result.Findings)
	}
}

func TestIntegration_OnlyWarnings_HasErrorsFalse(t *testing.T) {
	snap := &state.Snapshot{
		Resources: []state.Resource{
			{ID: "x", Type: "db", Provider: "", Fields: nil},
		},
	}
	result := lint.Run(snap)
	if result.HasErrors() {
		t.Fatal("expected no errors, only warnings")
	}
	for _, f := range result.Findings {
		if f.Severity != lint.SeverityWarning {
			t.Errorf("expected warning, got %s", f.Severity)
		}
	}
}
