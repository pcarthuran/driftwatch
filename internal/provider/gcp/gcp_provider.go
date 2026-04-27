package gcp

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// Fetcher defines the interface for fetching GCP resources.
type Fetcher interface {
	Fetch(ctx context.Context, projectID string) ([]provider.Resource, error)
}

// gcpProvider implements provider.Provider for GCP.
type gcpProvider struct {
	projectID string
	fetcher   Fetcher
}

// New creates a new GCP provider using the real fetcher.
func New(projectID string) (provider.Provider, error) {
	if projectID == "" {
		return nil, fmt.Errorf("gcp: project_id is required")
	}
	return &gcpProvider{
		projectID: projectID,
		fetcher:   &realFetcher{},
	}, nil
}

// NewWithFetcher creates a new GCP provider with a custom fetcher (for testing).
func NewWithFetcher(projectID string, fetcher Fetcher) (provider.Provider, error) {
	if projectID == "" {
		return nil, fmt.Errorf("gcp: project_id is required")
	}
	if fetcher == nil {
		return nil, fmt.Errorf("gcp: fetcher must not be nil")
	}
	return &gcpProvider{
		projectID: projectID,
		fetcher:   fetcher,
	}, nil
}

// Name returns the name of the provider.
func (p *gcpProvider) Name() string {
	return "gcp"
}

// Fetch retrieves all resources for the configured GCP project.
func (p *gcpProvider) Fetch(ctx context.Context) ([]provider.Resource, error) {
	resources, err := p.fetcher.Fetch(ctx, p.projectID)
	if err != nil {
		return nil, fmt.Errorf("gcp: fetch failed for project %q: %w", p.projectID, err)
	}
	return resources, nil
}
