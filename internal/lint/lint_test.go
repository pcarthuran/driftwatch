package lint_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/lint"
	"github.com/driftwatch/internal/state"
)

func baseSnapshot() *state.Snapshot {
	return &state.Snapshot{
		Resources: []state.Resource{
			{
				ID:       "res-1",
				Type:     "instance",
				Provider: "aws",
				Fields:   map[string]any{"region": "us-east-1"},
			},
		},
	}
}

func TestRun_NoDrift(t *testing.T) {
	snap := baseSnapshot()
	result := lint.Run(snap)
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(result.Findings))
	}
	if result.HasErrors() {
		t.Fatal("expected no errors")
	}
}

func TestRun_MissingID(t *testing.T) {
	snap := baseSnapshot()
	snap.Resources[0].ID = ""
	result := lint.Run(snap)
	if !result.HasErrors() {
		t.Fatal("expected error for missing id")
	}
	if !containsField(result, "id") {
		t.Error("expected finding for field 'id'")
	}
}

func TestRun_MissingType(t *testing.T) {
	snap := baseSnapshot()
	snap.Resources[0].Type = ""
	result := lint.Run(snap)
	if !result.HasErrors() {
		t.Fatal("expected error for missing type")
	}
}

func TestRun_MissingProvider_IsWarning(t *testing.T) {
	snap := baseSnapshot()
	snap.Resources[0].Provider = ""
	result := lint.Run(snap)
	if result.HasErrors() {
		t.Fatal("expected only warning, not error")
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(result.Findings))
	}
	if result.Findings[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity, got %s", result.Findings[0].Severity)
	}
}

func TestRun_EmptyFields_IsWarning(t *testing.T) {
	snap := baseSnapshot()
	snap.Resources[0].Fields = nil
	result := lint.Run(snap)
	if result.HasErrors() {
		t.Fatal("expected no errors")
	}
	if !containsField(result, "fields") {
		t.Error("expected warning for empty fields")
	}
}

func TestWrite_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	lint.Write(&buf, &lint.Result{})
	if !strings.Contains(buf.String(), "no issues") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWrite_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	r := &lint.Result{
		Findings: []lint.Finding{
			{ResourceID: "r1", Field: "id", Message: "empty", Severity: lint.SeverityError},
		},
	}
	lint.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "1 issue") {
		t.Errorf("expected issue count in output: %s", out)
	}
	if !strings.Contains(out, "[error]") {
		t.Errorf("expected severity in output: %s", out)
	}
}

func containsField(r *lint.Result, field string) bool {
	for _, f := range r.Findings {
		if f.Field == field {
			return true
		}
	}
	return false
}
