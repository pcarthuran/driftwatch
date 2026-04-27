package baseline_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/drift"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{
			ResourceID: "res-1",
			Status:     drift.StatusModified,
			Diffs:      []drift.FieldDiff{{Field: "size", Declared: "t2.micro", Live: "t2.large"}},
		},
		{
			ResourceID: "res-2",
			Status:     drift.StatusMatch,
		},
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/baselines"
	_, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSave_AndLatest(t *testing.T) {
	dir := t.TempDir()
	store, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	results := sampleResults()
	if err := store.Save("v1", results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := store.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if entry.Label != "v1" {
		t.Errorf("expected label 'v1', got %q", entry.Label)
	}
	if len(entry.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(entry.Results))
	}
	if entry.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	dir := t.TempDir()
	store, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	_ = store.Save("first", sampleResults())
	time.Sleep(2 * time.Millisecond)
	_ = store.Save("second", sampleResults())

	entry, err := store.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if entry.Label != "second" {
		t.Errorf("expected 'second', got %q", entry.Label)
	}
}

func TestLatest_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	store, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err = store.Latest()
	if err == nil {
		t.Error("expected error for empty store, got nil")
	}
}
