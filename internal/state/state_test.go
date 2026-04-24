package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/state"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoad_ValidJSON(t *testing.T) {
	content := `{"version":"1","resources":[{"id":"vpc-1","type":"vpc","provider":"aws","attributes":{"cidr":"10.0.0.0/16"}}]}`
	path := writeTempFile(t, "state.json", content)

	sf, err := state.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sf.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(sf.Resources))
	}
	if sf.Resources[0].ID != "vpc-1" {
		t.Errorf("expected id vpc-1, got %q", sf.Resources[0].ID)
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	content := "version: \"1\"\nresources:\n  - id: subnet-1\n    type: subnet\n    provider: aws\n    attributes:\n      az: us-east-1a\n"
	path := writeTempFile(t, "state.yaml", content)

	sf, err := state.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sf.Resources[0].ID != "subnet-1" {
		t.Errorf("expected id subnet-1, got %q", sf.Resources[0].ID)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := state.Load("/nonexistent/state.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_UnsupportedFormat(t *testing.T) {
	path := writeTempFile(t, "state.toml", "version = \"1\"")
	_, err := state.Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	content := `{"resources":[{"id":"vpc-1","type":"vpc"}]}`
	path := writeTempFile(t, "state.json", content)
	_, err := state.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing version")
	}
}

func TestResourceMap(t *testing.T) {
	content := `{"version":"1","resources":[{"id":"a","type":"vpc"},{"id":"b","type":"subnet"}]}`
	path := writeTempFile(t, "state.json", content)

	sf, err := state.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rm := sf.ResourceMap()
	if _, ok := rm["a"]; !ok {
		t.Error("expected key \"a\" in resource map")
	}
	if _, ok := rm["b"]; !ok {
		t.Error("expected key \"b\" in resource map")
	}
}
