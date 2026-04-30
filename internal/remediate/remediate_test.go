package remediate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/diff"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/remediate"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{
			ResourceID: "res-1",
			Provider:   "aws",
			Missing:    true,
		},
		{
			ResourceID: "res-2",
			Provider:   "gcp",
			Extra:      true,
		},
		{
			ResourceID: "res-3",
			Provider:   "azure",
			Diffs: []diff.FieldDiff{
				{Field: "size", Expected: "t2.micro", Actual: "t2.large"},
			},
		},
	}
}

func TestPlan_NoDrift(t *testing.T) {
	actions := remediate.Plan([]drift.Result{})
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestPlan_MissingResource(t *testing.T) {
	results := []drift.Result{{ResourceID: "r1", Provider: "aws", Missing: true}}
	actions := remediate.Plan(results)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Kind != "missing" {
		t.Errorf("expected kind 'missing', got %q", actions[0].Kind)
	}
	if !strings.Contains(actions[0].Suggestion, "r1") {
		t.Errorf("suggestion should reference resource ID, got: %s", actions[0].Suggestion)
	}
}

func TestPlan_ExtraResource(t *testing.T) {
	results := []drift.Result{{ResourceID: "r2", Provider: "gcp", Extra: true}}
	actions := remediate.Plan(results)
	if len(actions) != 1 || actions[0].Kind != "extra" {
		t.Errorf("expected 1 extra action, got %+v", actions)
	}
}

func TestPlan_ModifiedResource(t *testing.T) {
	results := []drift.Result{
		{
			ResourceID: "r3",
			Provider:   "azure",
			Diffs:      []diff.FieldDiff{{Field: "region", Expected: "us-east-1", Actual: "us-west-2"}},
		},
	}
	actions := remediate.Plan(results)
	if len(actions) != 1 || actions[0].Kind != "modified" {
		t.Fatalf("expected 1 modified action, got %+v", actions)
	}
	if !strings.Contains(actions[0].Suggestion, "region") {
		t.Errorf("suggestion should mention drifted field, got: %s", actions[0].Suggestion)
	}
}

func TestPlan_AllTypes(t *testing.T) {
	actions := remediate.Plan(sampleResults())
	if len(actions) != 3 {
		t.Errorf("expected 3 actions, got %d", len(actions))
	}
}

func TestWrite_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	remediate.Write(&buf, []remediate.Action{})
	if !strings.Contains(buf.String(), "No remediation") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWrite_WithActions(t *testing.T) {
	actions := remediate.Plan(sampleResults())
	var buf bytes.Buffer
	remediate.Write(&buf, actions)
	out := buf.String()
	if !strings.Contains(out, "Remediation Plan") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "MISSING") || !strings.Contains(out, "EXTRA") || !strings.Contains(out, "MODIFIED") {
		t.Errorf("expected all action kinds in output, got: %s", out)
	}
}
