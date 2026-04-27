package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a single drift detection run stored in history.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Provider  string         `json:"provider"`
	Results   []drift.Result `json:"results"`
	DriftCount int           `json:"drift_count"`
}

// Store manages persisting and retrieving drift history entries.
type Store struct {
	Dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir %q: %w", dir, err)
	}
	return &Store{Dir: dir}, nil
}

// Save writes a new history entry to disk as a timestamped JSON file.
func (s *Store) Save(entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	filename := fmt.Sprintf("%s_%s.json",
		entry.Timestamp.Format("20060102T150405Z"),
		sanitize(entry.Provider),
	)
	path := filepath.Join(s.Dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("history: encode entry: %w", err)
	}
	return nil
}

// List returns all history entries sorted by timestamp descending.
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.Dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []Entry
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("history: read %q: %w", path, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: decode %q: %w", path, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	return entries, nil
}

func sanitize(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return string(out)
}
