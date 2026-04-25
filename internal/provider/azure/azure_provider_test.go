package azure_test

import (
	"errors"
	"testing"

	"github.com/user/driftwatch/internal/provider/azure"
	"github.com/user/driftwatch/internal/state"
)

type mockFetcher struct {
	resources []state.Resource
	err       error
}

func (m *mockFetcher) Fetch() ([]state.Resource, error) {
	return m.resources, m.err
}

func sampleResources() []state.Resource {
	return []state.Resource{
		{
			ID:   "vm-001",
			Type: "azure_virtual_machine",
			Fields: map[string]interface{}{
				"location": "eastus",
				"size":     "Standard_D2s_v3",
			},
		},
		{
			ID:   "rg-001",
			Type: "azure_resource_group",
			Fields: map[string]interface{}{
				"location": "westus",
			},
		},
	}
}

func TestAzureProvider_Name(t *testing.T) {
	p := azure.NewWithFetcher("sub-123", &mockFetcher{})
	if p.Name() != "azure" {
		t.Errorf("expected name 'azure', got '%s'", p.Name())
	}
}

func TestAzureProvider_Fetch_Success(t *testing.T) {
	resources := sampleResources()
	p := azure.NewWithFetcher("sub-123", &mockFetcher{resources: resources})

	got, err := p.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(resources) {
		t.Errorf("expected %d resources, got %d", len(resources), len(got))
	}
}

func TestAzureProvider_Fetch_Error(t *testing.T) {
	p := azure.NewWithFetcher("sub-123", &mockFetcher{err: errors.New("azure api error")})

	_, err := p.Fetch()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAzureProvider_SubscriptionID(t *testing.T) {
	p := azure.NewWithFetcher("sub-abc", &mockFetcher{})
	if p.SubscriptionID() != "sub-abc" {
		t.Errorf("expected subscription ID 'sub-abc', got '%s'", p.SubscriptionID())
	}
}

func TestAzureProvider_Fetch_Empty(t *testing.T) {
	p := azure.NewWithFetcher("sub-123", &mockFetcher{resources: []state.Resource{}})

	got, err := p.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 resources, got %d", len(got))
	}
}
