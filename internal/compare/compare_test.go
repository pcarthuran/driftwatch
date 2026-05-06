package compare_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/compare"
	"github.com/driftwatch/internal/snapshot"
	"github.com/driftwatch/internal/state"
)

func writeTempSnapshot(t *testing.T, resources []state.Resource) string {
	t.Helper()
	snap := snapshot.Snapshot{Resources: resources}
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

var baseResources = []state.Resource{
	{ID: "res-1", Type: "vm", Provider: "aws", Fields: map[string]any{"size": "t2.micro"}},
	{ID: "res-2", Type: "bucket", Provider: "aws", Fields: map[string]any{"region": "us-east-1"}},
}

func TestCompare_NoDrift(t *testing.T) {
	a := writeTempSnapshot(t, baseResources)
	b := writeTempSnapshot(t, baseResources)

	result, err := compare.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(result.Diffs))
	}
}

func TestCompare_DetectsDrift(t *testing.T) {
	modified := []state.Resource{
		{ID: "res-1", Type: "vm", Provider: "aws", Fields: map[string]any{"size": "t3.large"}},
		{ID: "res-2", Type: "bucket", Provider: "aws", Fields: map[string]any{"region": "us-east-1"}},
	}
	a := writeTempSnapshot(t, baseResources)
	b := writeTempSnapshot(t, modified)

	result, err := compare.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(result.Diffs))
	}
}

func TestCompare_MissingBaselineFile(t *testing.T) {
	current := writeTempSnapshot(t, baseResources)
	_, err := compare.Compare(filepath.Join(t.TempDir(), "missing.json"), current)
	if err == nil {
		t.Error("expected error for missing baseline, got nil")
	}
}

func TestWrite_NoDrift(t *testing.T) {
	r := &compare.Result{BaselineLabel: "a.json", CurrentLabel: "b.json", Diffs: nil}
	var buf bytes.Buffer
	compare.Write(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("No drift")) {
		t.Errorf("expected 'No drift' in output, got: %s", buf.String())
	}
}

func TestWrite_WithDrift(t *testing.T) {
	a := writeTempSnapshot(t, baseResources)
	modified := []state.Resource{
		{ID: "res-1", Type: "vm", Provider: "aws", Fields: map[string]any{"size": "t3.large"}},
	}
	b := writeTempSnapshot(t, modified)

	result, _ := compare.Compare(a, b)
	var buf bytes.Buffer
	compare.Write(&buf, result)
	if !bytes.Contains(buf.Bytes(), []byte("res-1")) {
		t.Errorf("expected resource ID in output, got: %s", buf.String())
	}
}
