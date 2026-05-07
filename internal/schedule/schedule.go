package schedule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a scheduled drift-check job.
type Entry struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Provider  string        `json:"provider"`
	StateFile string        `json:"state_file"`
	Interval  time.Duration `json:"interval_ns"`
	LastRun   time.Time     `json:"last_run,omitempty"`
	Enabled   bool          `json:"enabled"`
}

// Store manages persisted schedule entries.
type Store struct {
	dir string
}

// NewStore creates (or opens) a schedule store at the given directory.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("schedule: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save persists an entry to disk, overwriting any existing entry with the same ID.
func (s *Store) Save(e Entry) error {
	if e.ID == "" {
		return fmt.Errorf("schedule: entry ID must not be empty")
	}
	path := filepath.Join(s.dir, e.ID+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("schedule: save entry: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}

// List returns all schedule entries found in the store directory.
func (s *Store) List() ([]Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("schedule: list entries: %w", err)
	}
	var entries []Entry
	for _, path := range matches {
		var e Entry
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("schedule: read %s: %w", path, err)
		}
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("schedule: decode %s: %w", path, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Get returns a single schedule entry by ID.
// It returns an error wrapping os.ErrNotExist if no entry with that ID is found.
func (s *Store) Get(id string) (Entry, error) {
	path := filepath.Join(s.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, fmt.Errorf("schedule: get %s: %w", id, err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, fmt.Errorf("schedule: decode %s: %w", id, err)
	}
	return e, nil
}

// Delete removes a schedule entry by ID.
func (s *Store) Delete(id string) error {
	path := filepath.Join(s.dir, id+".json")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("schedule: delete %s: %w", id, err)
	}
	return nil
}
