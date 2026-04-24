package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the driftwatch configuration.
type Config struct {
	StateFile  string            `yaml:"state_file"`
	Provider   string            `yaml:"provider"`
	Ignore     []string          `yaml:"ignore"`
	Labels     map[string]string `yaml:"labels"`
	OutputFormat string          `yaml:"output_format"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		StateFile:    "drift.state.yaml",
		Provider:     "local",
		Ignore:       []string{},
		Labels:       map[string]string{},
		OutputFormat: "text",
	}
}

// Load reads a YAML config file from the given path and returns a Config.
// If the path is empty, it searches for a default config file in the current
// directory and the user's home directory.
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = findDefaultConfig()
		if err != nil {
			return DefaultConfig(), nil
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: opening %q: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: parsing %q: %w", path, err)
	}

	return cfg, nil
}

// findDefaultConfig searches well-known locations for a driftwatch config file.
func findDefaultConfig() (string, error) {
	candidates := []string{
		".driftwatch.yaml",
		".driftwatch.yml",
	}

	home, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".driftwatch.yaml"),
			filepath.Join(home, ".driftwatch.yml"),
		)
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}

	return "", fmt.Errorf("config: no default config file found")
}
