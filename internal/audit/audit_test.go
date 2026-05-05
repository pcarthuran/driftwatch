package audit_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/audit"
)

func sampleEntry(id string, ts time.Time) audit.Entry {
	return audit.Entry{
		ID:        id,
		Timestamp: ts,
		Provider:  "aws",
		Total:     10,
		Drifted:   2,
		Missing:   1,
		Extra:     1,
		User:      "ci-bot",
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/audit/logs"
	_, err := audit.NewStore(subdir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	store, _ := audit.NewStore(dir)
	e := sampleEntry("run-001", time.Now())
	if err := store.Save(e); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "run-001" {
		t.Errorf("expected ID run-001, got %s", entries[0].ID)
	}
}

func TestSave_EmptyID_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	store, _ := audit.NewStore(dir)
	e := sampleEntry("", time.Now())
	if err := store.Save(e); err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}
}

func TestList_SortedDescending(t *testing.T) {
	dir := t.TempDir()
	store, _ := audit.NewStore(dir)
	now := time.Now().UTC()
	_ = store.Save(sampleEntry("run-001", now.Add(-2*time.Hour)))
	_ = store.Save(sampleEntry("run-002", now.Add(-1*time.Hour)))
	_ = store.Save(sampleEntry("run-003", now))
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].ID != "run-003" {
		t.Errorf("expected newest first, got %s", entries[0].ID)
	}
	if entries[2].ID != "run-001" {
		t.Errorf("expected oldest last, got %s", entries[2].ID)
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	store, _ := audit.NewStore(dir)
	entries, err := store.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
