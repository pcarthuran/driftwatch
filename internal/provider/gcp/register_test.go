package gcp_test

import (
	"testing"

	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/provider/gcp"
)

func TestRegister_Success(t *testing.T) {
	registry := provider.NewRegistry()
	opts := map[string]string{"project_id": "my-gcp-project"}

	if err := gcp.Register(registry, opts); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	p, err := registry.Get("gcp")
	if err != nil {
		t.Fatalf("expected provider to be registered, got: %v", err)
	}
	if p.Name() != "gcp" {
		t.Errorf("expected provider name 'gcp', got %q", p.Name())
	}
}

func TestRegister_MissingProjectID(t *testing.T) {
	registry := provider.NewRegistry()

	tests := []struct {
		name string
		opts map[string]string
	}{
		{"nil opts", map[string]string{}},
		{"empty project_id", map[string]string{"project_id": ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := gcp.Register(registry, tt.opts); err == nil {
				t.Error("expected error for missing project_id, got nil")
			}
		})
	}
}

func TestRegister_Duplicate(t *testing.T) {
	registry := provider.NewRegistry()
	opts := map[string]string{"project_id": "my-gcp-project"}

	if err := gcp.Register(registry, opts); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	if err := gcp.Register(registry, opts); err == nil {
		t.Error("expected error on duplicate registration, got nil")
	}
}
