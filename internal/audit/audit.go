// Package audit provides functionality for recording and retrieving
// drift detection audit log entries.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a single audit log record for a drift detection run.
type Entry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Provider  string    `json:"provider"`
	Total     int       `json:"total_resources"`
	Drifted   int       `json:"drifted"`
	Missing   int       `json:"missing"`
	Extra     int       `json:"extra"`
	User      string    `json:"user,omitempty"`
}

// Store manages persistence of audit log entries on disk.
type Store struct {
	dir string
}

// NewStore creates a new Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("audit: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes an audit entry to disk as a JSON file named by entry ID and timestamp.
func (s *Store) Save(e Entry) error {
	if e.ID == "" {
		return fmt.Errorf("audit: entry ID must not be empty")
	}
	filename := fmt.Sprintf("%s_%s.json", e.Timestamp.UTC().Format("20060102T150405Z"), e.ID)
	path := filepath.Join(s.dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("audit: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// List returns all audit entries sorted by timestamp descending.
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("audit: glob: %w", err)
	}
	var entries []Entry
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("audit: read %s: %w", path, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("audit: unmarshal %s: %w", path, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	return entries, nil
}
