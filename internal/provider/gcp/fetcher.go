package gcp

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// realFetcher is the production GCP resource fetcher.
type realFetcher struct {
	// projectID is the GCP project to fetch resources from.
	projectID string
}

// NewFetcher creates a new GCP resource fetcher for the given project.
func NewFetcher(projectID string) provider.Fetcher {
	return &realFetcher{projectID: projectID}
}

// Fetch retrieves resources from GCP for the given project.
// In a real implementation this would call the GCP API (e.g., Asset Inventory).
// The context can be used to cancel long-running API calls.
func (f *realFetcher) Fetch(ctx context.Context, projectID string) ([]provider.Resource, error) {
	if projectID == "" {
		return nil, fmt.Errorf("gcp fetcher: projectID must not be empty")
	}
	// TODO: integrate with GCP Asset Inventory or specific service APIs.
	return nil, fmt.Errorf("gcp real fetcher not yet implemented for project %s", projectID)
}
