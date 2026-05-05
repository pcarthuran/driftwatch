package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Snapshot represents a point-in-time capture of infrastructure state.
type Snapshot struct {
	CapturedAt time.Time                `json:"captured_at" yaml:"captured_at"`
	Source     string                   `json:"source" yaml:"source"`
	Resources  []map[string]interface{} `json:"resources" yaml:"resources"`
}

// Save writes the snapshot to the given file path.
// The format is determined by the file extension (.json or .yaml/.yml).
func Save(snap *Snapshot, path string) error {
	snap.CapturedAt = time.Now().UTC()

	var data []byte
	var err error

	switch filepath.Ext(path) {
	case ".json":
		data, err = json.MarshalIndent(snap, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(snap)
	default:
		return fmt.Errorf("unsupported snapshot format: %s", filepath.Ext(path))
	}

	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snap Snapshot

	switch filepath.Ext(path) {
	case ".json":
		if err := json.Unmarshal(data, &snap); err != nil {
			return nil, fmt.Errorf("failed to parse snapshot JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &snap); err != nil {
			return nil, fmt.Errorf("failed to parse snapshot YAML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported snapshot format: %s", filepath.Ext(path))
	}

	return &snap, nil
}

// ResourceCount returns the number of resources captured in the snapshot.
func (s *Snapshot) ResourceCount() int {
	return len(s.Resources)
}

// Age returns the duration elapsed since the snapshot was captured.
func (s *Snapshot) Age() time.Duration {
	return time.Since(s.CapturedAt)
}
