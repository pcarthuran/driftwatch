package azure

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/provider"
)

const ProviderName = "azure"

// Fetcher defines the interface for fetching Azure resources.
type Fetcher interface {
	FetchResources(ctx context.Context, subscriptionID string) ([]provider.Resource, error)
}

// AzureProvider fetches live resource state from Azure.
type AzureProvider struct {
	subscriptionID string
	fetcher        Fetcher
}

// New creates an AzureProvider using the real Azure fetcher.
func New(subscriptionID string) (*AzureProvider, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("azure: subscription_id is required")
	}
	return &AzureProvider{
		subscriptionID: subscriptionID,
		fetcher:        &realFetcher{},
	}, nil
}

// NewWithFetcher creates an AzureProvider with a custom fetcher (useful for testing).
func NewWithFetcher(subscriptionID string, fetcher Fetcher) (*AzureProvider, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("azure: subscription_id is required")
	}
	return &AzureProvider{
		subscriptionID: subscriptionID,
		fetcher:        fetcher,
	}, nil
}

// Name returns the provider identifier.
func (p *AzureProvider) Name() string {
	return ProviderName
}

// Fetch retrieves live resources from Azure.
func (p *AzureProvider) Fetch(ctx context.Context) ([]provider.Resource, error) {
	resources, err := p.fetcher.FetchResources(ctx, p.subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("azure: failed to fetch resources: %w", err)
	}
	return resources, nil
}
