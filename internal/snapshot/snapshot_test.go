package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/snapshot"
)

func sampleSnapshot() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Source: "test-env",
		Resources: []map[string]interface{}{
			{"id": "res-1", "type": "vm", "status": "running"},
			{"id": "res-2", "type": "bucket", "region": "us-east-1"},
		},
	}
}

func TestSaveAndLoad_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	snap := sampleSnapshot()
	if err := snapshot.Save(snap, path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Source != snap.Source {
		t.Errorf("expected source %q, got %q", snap.Source, loaded.Source)
	}
	if len(loaded.Resources) != len(snap.Resources) {
		t.Errorf("expected %d resources, got %d", len(snap.Resources), len(loaded.Resources))
	}
	if loaded.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestSaveAndLoad_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.yaml")

	snap := sampleSnapshot()
	if err := snapshot.Save(snap, path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Source != snap.Source {
		t.Errorf("expected source %q, got %q", snap.Source, loaded.Source)
	}
}

func TestSave_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.toml")

	if err := snapshot.Save(sampleSnapshot(), path); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.toml")
	_ = os.WriteFile(path, []byte("source = 'test'"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
