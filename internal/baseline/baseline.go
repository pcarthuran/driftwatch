package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a saved baseline of drift detection results.
type Entry struct {
	CreatedAt time.Time           `json:"created_at"`
	Label     string              `json:"label,omitempty"`
	Results   []drift.Result      `json:"results"`
}

// Store manages baseline entries on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("baseline: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes results as a new baseline file under the store directory.
func (s *Store) Save(label string, results []drift.Result) error {
	entry := Entry{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	filename := fmt.Sprintf("%d.json", entry.CreatedAt.UnixNano())
	path := filepath.Join(s.dir, filename)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write file: %w", err)
	}
	return nil
}

// Latest returns the most recently saved baseline entry, or an error if none exist.
func (s *Store) Latest() (*Entry, error) {
	entries, err := s.list()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("baseline: no entries found in %s", s.dir)
	}
	return entries[len(entries)-1], nil
}

// list returns all baseline entries sorted ascending by filename (timestamp).
func (s *Store) list() ([]*Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("baseline: glob: %w", err)
	}
	var entries []*Entry
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("baseline: read %s: %w", path, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("baseline: unmarshal %s: %w", path, err)
		}
		entries = append(entries, &e)
	}
	return entries, nil
}
