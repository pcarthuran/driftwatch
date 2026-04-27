package schedule_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/schedule"
)

func sampleEntry(id string) schedule.Entry {
	return schedule.Entry{
		ID:        id,
		Name:      "nightly-check",
		Provider:  "aws",
		StateFile: "infra/main.yaml",
		Interval:  24 * time.Hour,
		Enabled:   true,
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/schedules"
	_, err := schedule.NewStore(subdir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSave_AndList(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	e := sampleEntry("job-001")
	if err := store.Save(e); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "job-001" {
		t.Errorf("expected ID job-001, got %s", entries[0].ID)
	}
}

func TestSave_EmptyID_ReturnsError(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	e := sampleEntry("")
	if err := store.Save(e); err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	e := sampleEntry("job-002")
	_ = store.Save(e)
	if err := store.Delete("job-002"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	entries, _ := store.List()
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after delete, got %d", len(entries))
	}
}

func TestDelete_MissingID_ReturnsError(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	if err := store.Delete("nonexistent"); err == nil {
		t.Error("expected error when deleting missing entry")
	}
}

func TestList_EmptyDir(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	entries, err := store.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
