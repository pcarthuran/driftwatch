package gcp

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// realFetcher is the production GCP resource fetcher.
type realFetcher struct{}

// Fetch retrieves resources from GCP for the given project.
// In a real implementation this would call the GCP API (e.g., Asset Inventory).
func (f *realFetcher) Fetch(ctx context.Context, projectID string) ([]provider.Resource, error) {
	// TODO: integrate with GCP Asset Inventory or specific service APIs.
	return nil, fmt.Errorf("gcp real fetcher not yet implemented for project %s", projectID)
}
