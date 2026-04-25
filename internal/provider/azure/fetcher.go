package azure

import (
	"fmt"

	"github.com/user/driftwatch/internal/state"
)

// Fetcher defines the interface for retrieving Azure resources.
type Fetcher interface {
	Fetch() ([]state.Resource, error)
}

// liveFetcher retrieves resources from the Azure API using the subscription ID.
type liveFetcher struct {
	subscriptionID string
}

// newLiveFetcher creates a new liveFetcher for the given Azure subscription.
func newLiveFetcher(subscriptionID string) *liveFetcher {
	return &liveFetcher{subscriptionID: subscriptionID}
}

// Fetch retrieves live Azure resources for the configured subscription.
// NOTE: This is a stub — replace with real Azure SDK calls (e.g., armresources).
func (f *liveFetcher) Fetch() ([]state.Resource, error) {
	if f.subscriptionID == "" {
		return nil, fmt.Errorf("azure: subscription ID must not be empty")
	}
	// TODO: integrate github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources
	return []state.Resource{}, nil
}
