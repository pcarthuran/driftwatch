package gcp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/provider/gcp"
)

// mockFetcher is a test double for gcp.Fetcher.
type mockFetcher struct {
	resources []provider.Resource
	err       error
}

func (m *mockFetcher) Fetch(_ context.Context, _ string) ([]provider.Resource, error) {
	return m.resources, m.err
}

func sampleResources() []provider.Resource {
	return []provider.Resource{
		{ID: "instance-1", Type: "compute.Instance", Fields: map[string]interface{}{"zone": "us-central1-a", "status": "RUNNING"}},
		{ID: "bucket-1", Type: "storage.Bucket", Fields: map[string]interface{}{"location": "US", "storageClass": "STANDARD"}},
	}
}

func TestGCPProvider_Name(t *testing.T) {
	p, err := gcp.NewWithFetcher("my-project", &mockFetcher{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Name(); got != "gcp" {
		t.Errorf("expected name 'gcp', got %q", got)
	}
}

func TestGCPProvider_Fetch_Success(t *testing.T) {
	resources := sampleResources()
	p, err := gcp.NewWithFetcher("my-project", &mockFetcher{resources: resources})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}
	if len(got) != len(resources) {
		t.Errorf("expected %d resources, got %d", len(resources), len(got))
	}
}

func TestGCPProvider_Fetch_Empty(t *testing.T) {
	p, err := gcp.NewWithFetcher("my-project", &mockFetcher{resources: []provider.Resource{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 resources, got %d", len(got))
	}
}

func TestGCPProvider_Fetch_Error(t *testing.T) {
	fetchErr := errors.New("gcp api unavailable")
	p, err := gcp.NewWithFetcher("my-project", &mockFetcher{err: fetchErr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGCPProvider_MissingProjectID(t *testing.T) {
	_, err := gcp.NewWithFetcher("", &mockFetcher{})
	if err == nil {
		t.Fatal("expected error for empty project_id, got nil")
	}
}

func TestNew_MissingProjectID(t *testing.T) {
	_, err := gcp.New("")
	if err == nil {
		t.Fatal("expected error for empty project_id, got nil")
	}
}
