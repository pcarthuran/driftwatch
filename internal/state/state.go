package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Resource represents a single declared infrastructure resource.
type Resource struct {
	ID         string            `json:"id" yaml:"id"`
	Type       string            `json:"type" yaml:"type"`
	Provider   string            `json:"provider" yaml:"provider"`
	Attributes map[string]string `json:"attributes" yaml:"attributes"`
}

// StateFile represents the declared state loaded from a file.
type StateFile struct {
	Version   string     `json:"version" yaml:"version"`
	Resources []Resource `json:"resources" yaml:"resources"`
}

// Load reads a state file from disk. Supports .json and .yaml/.yml formats.
func Load(path string) (*StateFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", path, err)
	}

	var sf StateFile
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &sf); err != nil {
			return nil, fmt.Errorf("parsing JSON state file: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &sf); err != nil {
			return nil, fmt.Errorf("parsing YAML state file: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported state file format %q (use .json or .yaml)", ext)
	}

	if err := sf.validate(); err != nil {
		return nil, fmt.Errorf("invalid state file: %w", err)
	}

	return &sf, nil
}

// validate performs basic sanity checks on the loaded state.
func (sf *StateFile) validate() error {
	if sf.Version == "" {
		return fmt.Errorf("missing required field \"version\"")
	}
	for i, r := range sf.Resources {
		if r.ID == "" {
			return fmt.Errorf("resource[%d] missing required field \"id\"", i)
		}
		if r.Type == "" {
			return fmt.Errorf("resource %q missing required field \"type\"", r.ID)
		}
	}
	return nil
}

// ResourceMap returns resources indexed by their ID for quick lookup.
func (sf *StateFile) ResourceMap() map[string]Resource {
	m := make(map[string]Resource, len(sf.Resources))
	for _, r := range sf.Resources {
		m[r.ID] = r
	}
	return m
}
