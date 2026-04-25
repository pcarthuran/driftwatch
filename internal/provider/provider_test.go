package provider_test

import (
	"errors"
	"testing"

	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/state"
)

// mockProvider is a test double implementing Provider.
type mockProvider struct {
	name      string
	resources []state.Resource
	err       error
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) FetchResources() ([]state.Resource, error) {
	return m.resources, m.err
}

func sampleResources() []state.Resource {
	return []state.Resource{
		{ID: "res-1", Type: "vm", Fields: map[string]interface{}{"size": "small"}},
		{ID: "res-2", Type: "db", Fields: map[string]interface{}{"engine": "postgres"}},
	}
}

func TestRegister_Success(t *testing.T) {
	reg := provider.NewRegistry()
	err := reg.Register(&mockProvider{name: "aws"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	reg := provider.NewRegistry()
	_ = reg.Register(&mockProvider{name: "aws"})
	err := reg.Register(&mockProvider{name: "aws"})
	if err == nil {
		t.Fatal("expected duplicate registration error")
	}
}

func TestRegister_EmptyName(t *testing.T) {
	reg := provider.NewRegistry()
	err := reg.Register(&mockProvider{name: ""})
	if err == nil {
		t.Fatal("expected error for empty provider name")
	}
}

func TestGet_Found(t *testing.T) {
	reg := provider.NewRegistry()
	_ = reg.Register(&mockProvider{name: "gcp"})
	p, err := reg.Get("gcp")
	if err != nil {
		t.Fatalf("expected provider, got error: %v", err)
	}
	if p.Name() != "gcp" {
		t.Errorf("expected name gcp, got %s", p.Name())
	}
}

func TestGet_NotFound(t *testing.T) {
	reg := provider.NewRegistry()
	_, err := reg.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing provider")
	}
}

func TestFetchAll_MergesResources(t *testing.T) {
	reg := provider.NewRegistry()
	_ = reg.Register(&mockProvider{name: "aws", resources: sampleResources()[:1]})
	_ = reg.Register(&mockProvider{name: "gcp", resources: sampleResources()[1:]})

	resources, err := reg.FetchAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(resources))
	}
}

func TestFetchAll_PropagatesError(t *testing.T) {
	reg := provider.NewRegistry()
	_ = reg.Register(&mockProvider{name: "aws", err: errors.New("connection refused")})

	_, err := reg.FetchAll()
	if err == nil {
		t.Fatal("expected fetch error to propagate")
	}
}
