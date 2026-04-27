package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/history"
)

func sampleEntry(provider string, driftCount int) history.Entry {
	return history.Entry{
		Timestamp:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Provider:   provider,
		DriftCount: driftCount,
		Results: []drift.Result{
			{ResourceID: "res-1", Status: "modified"},
		},
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/history"
	_, err := history.NewStore(subdir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	entry := sampleEntry("aws", 1)
	if err := store.Save(entry); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Provider != "aws" {
		t.Errorf("expected provider aws, got %s", entries[0].Provider)
	}
	if entries[0].DriftCount != 1 {
		t.Errorf("expected drift_count 1, got %d", entries[0].DriftCount)
	}
}

func TestList_SortedDescending(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	e1 := sampleEntry("aws", 0)
	e1.Timestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e2 := sampleEntry("gcp", 2)
	e2.Timestamp = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	if err := store.Save(e1); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(e2); err != nil {
		t.Fatal(err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Provider != "gcp" {
		t.Errorf("expected gcp first (newest), got %s", entries[0].Provider)
	}
}

func TestList_EmptyDir(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
