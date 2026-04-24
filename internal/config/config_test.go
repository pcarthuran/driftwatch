package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.StateFile != "drift.state.yaml" {
		t.Errorf("expected StateFile %q, got %q", "drift.state.yaml", cfg.StateFile)
	}
	if cfg.Provider != "local" {
		t.Errorf("expected Provider %q, got %q", "local", cfg.Provider)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected OutputFormat %q, got %q", "text", cfg.OutputFormat)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "test.yaml")

	content := []byte(`
state_file: custom.state.yaml
provider: aws
output_format: json
ignore:
  - "*.tmp"
labels:
  env: staging
`)
	if err := os.WriteFile(cfgPath, content, 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.StateFile != "custom.state.yaml" {
		t.Errorf("expected StateFile %q, got %q", "custom.state.yaml", cfg.StateFile)
	}
	if cfg.Provider != "aws" {
		t.Errorf("expected Provider %q, got %q", "aws", cfg.Provider)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected OutputFormat %q, got %q", "json", cfg.OutputFormat)
	}
	if len(cfg.Ignore) != 1 || cfg.Ignore[0] != "*.tmp" {
		t.Errorf("unexpected Ignore list: %v", cfg.Ignore)
	}
	if cfg.Labels["env"] != "staging" {
		t.Errorf("expected label env=staging, got %v", cfg.Labels)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_NoPathFallsBackToDefaults(t *testing.T) {
	// Ensure no default config files exist in cwd by running from a temp dir.
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir) //nolint:errcheck

	t.Chdir(t.TempDir())

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}
	if cfg.Provider != "local" {
		t.Errorf("expected default provider %q, got %q", "local", cfg.Provider)
	}
}
