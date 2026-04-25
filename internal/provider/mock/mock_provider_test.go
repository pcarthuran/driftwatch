package mock_test

import (
	"testing"

	"github.com/user/driftwatch/internal/provider"
	"github.com/user/driftwatch/internal/provider/mock"
)

func sampleResources() []provider.Resource {
	return []provider.Resource{
		{
			ID:   "res-001",
			Type: "vm",
			Fields: map[string]interface{}{
				"region": "us-east-1",
				"size":   "t2.micro",
			},
		},
		{
			ID:   "res-002",
			Type: "bucket",
			Fields: map[string]interface{}{
				"region": "eu-west-1",
			},
		},
	}
}

func TestMockProvider_Name(t *testing.T) {
	p := mock.NewProvider("test-provider", nil)
	if p.Name() != "test-provider" {
		t.Errorf("expected name %q, got %q", "test-provider", p.Name())
	}
}

func TestMockProvider_Fetch_Success(t *testing.T) {
	resources := sampleResources()
	p := mock.NewProvider("aws-mock", resources)

	got, err := p.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(resources) {
		t.Errorf("expected %d resources, got %d", len(resources), len(got))
	}
	if got[0].ID != "res-001" {
		t.Errorf("expected first resource ID %q, got %q", "res-001", got[0].ID)
	}
}

func TestMockProvider_Fetch_Failure(t *testing.T) {
	p := mock.NewFailingProvider("broken-provider")

	_, err := p.Fetch()
	if err == nil {
		t.Fatal("expected error from failing provider, got nil")
	}
}

func TestMockProvider_SetResources(t *testing.T) {
	p := mock.NewProvider("dynamic", nil)

	got, err := p.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 resources initially, got %d", len(got))
	}

	p.SetResources(sampleResources())

	got, err = p.Fetch()
	if err != nil {
		t.Fatalf("unexpected error after SetResources: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 resources after SetResources, got %d", len(got))
	}
}
