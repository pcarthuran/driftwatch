package schedule_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/schedule"
)

func TestRoundTrip_MultipleEntries(t *testing.T) {
	store, err := schedule.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	ids := []string{"alpha", "beta", "gamma"}
	for _, id := range ids {
		e := schedule.Entry{
			ID:        id,
			Name:      id + "-job",
			Provider:  "gcp",
			StateFile: "state/" + id + ".yaml",
			Interval:  6 * time.Hour,
			Enabled:   true,
		}
		if err := store.Save(e); err != nil {
			t.Fatalf("Save(%s): %v", id, err)
		}
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != len(ids) {
		t.Fatalf("expected %d entries, got %d", len(ids), len(entries))
	}
}

func TestOverwrite_SameID(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	e := sampleEntry("overwrite-test")
	_ = store.Save(e)

	e.Name = "updated-name"
	e.Interval = 12 * time.Hour
	if err := store.Save(e); err != nil {
		t.Fatalf("second Save failed: %v", err)
	}

	entries, _ := store.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after overwrite, got %d", len(entries))
	}
	if entries[0].Name != "updated-name" {
		t.Errorf("expected updated name, got %s", entries[0].Name)
	}
	if entries[0].Interval != 12*time.Hour {
		t.Errorf("expected 12h interval, got %v", entries[0].Interval)
	}
}

func TestDeleteThenList(t *testing.T) {
	store, _ := schedule.NewStore(t.TempDir())
	_ = store.Save(sampleEntry("keep"))
	_ = store.Save(sampleEntry("remove"))
	_ = store.Delete("remove")

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "keep" {
		t.Errorf("expected 'keep', got %s", entries[0].ID)
	}
}
