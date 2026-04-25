package aws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/driftwatch/internal/provider"
	"github.com/user/driftwatch/internal/provider/aws"
)

// mockFetcher is a test double for ResourceFetcher.
type mockFetcher struct {
	resources []provider.Resource
	err       error
}

func (m *mockFetcher) FetchEC2Instances(_ context.Context, _ string) ([]provider.Resource, error) {
	return m.resources, m.err
}

func sampleResources() []provider.Resource {
	return []provider.Resource{
		{ID: "i-001", Type: "ec2_instance", Attributes: map[string]interface{}{"env": "prod", "state": "running"}},
		{ID: "i-002", Type: "ec2_instance", Attributes: map[string]interface{}{"env": "staging", "state": "stopped"}},
	}
}

func TestAWSProvider_Name(t *testing.T) {
	p := aws.New(aws.Config{Region: "us-east-1"})
	if p.Name() != "aws" {
		t.Errorf("expected name 'aws', got %q", p.Name())
	}
}

func TestAWSProvider_Fetch_Success(t *testing.T) {
	fetcher := &mockFetcher{resources: sampleResources()}
	p := aws.NewWithFetcher(aws.Config{Region: "us-east-1"}, fetcher)

	res, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 resources, got %d", len(res))
	}
}

func TestAWSProvider_Fetch_Error(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("api failure")}
	p := aws.NewWithFetcher(aws.Config{Region: "us-east-1"}, fetcher)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAWSProvider_Fetch_TagFilter(t *testing.T) {
	fetcher := &mockFetcher{resources: sampleResources()}
	cfg := aws.Config{
		Region:    "us-east-1",
		TagFilter: map[string]string{"env": "prod"},
	}
	p := aws.NewWithFetcher(cfg, fetcher)

	res, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("expected 1 resource after tag filter, got %d", len(res))
	}
	if res[0].ID != "i-001" {
		t.Errorf("expected resource i-001, got %s", res[0].ID)
	}
}

func TestAWSProvider_Fetch_TagFilter_NoMatch(t *testing.T) {
	fetcher := &mockFetcher{resources: sampleResources()}
	cfg := aws.Config{
		Region:    "us-east-1",
		TagFilter: map[string]string{"env": "dev"},
	}
	p := aws.NewWithFetcher(cfg, fetcher)

	res, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 resources, got %d", len(res))
	}
}
