package gcp

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/provider"
)

const providerName = "gcp"

// ResourceFetcher defines the interface for fetching GCP resources.
type ResourceFetcher interface {
	Fetch(ctx context.Context, projectID string) ([]provider.Resource, error)
}

// Provider implements provider.Provider for GCP.
type Provider struct {
	projectID string
	fetcher   ResourceFetcher
}

// New creates a GCP provider using the real fetcher.
func New(projectID string) (*Provider, error) {
	if projectID == "" {
		return nil, fmt.Errorf("gcp: project_id is required")
	}
	return &Provider{
		projectID: projectID,
		fetcher:   &realFetcher{},
	}, nil
}

// NewWithFetcher creates a GCP provider with a custom fetcher (for testing).
func NewWithFetcher(projectID string, fetcher ResourceFetcher) (*Provider, error) {
	if projectID == "" {
		return nil, fmt.Errorf("gcp: project_id is required")
	}
	return &Provider{
		projectID: projectID,
		fetcher:   fetcher,
	}, nil
}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return providerName
}

// Fetch retrieves live resources from GCP.
func (p *Provider) Fetch(ctx context.Context) ([]provider.Resource, error) {
	resources, err := p.fetcher.Fetch(ctx, p.projectID)
	if err != nil {
		return nil, fmt.Errorf("gcp: fetch failed for project %q: %w", p.projectID, err)
	}
	return resources, nil
}
