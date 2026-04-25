package mock

import (
	"fmt"

	"github.com/user/driftwatch/internal/provider"
)

// Provider is a mock infrastructure provider for testing and dry-run scenarios.
type Provider struct {
	resources []provider.Resource
	failFetch bool
	name      string
}

// NewProvider creates a new mock provider with the given resources.
func NewProvider(name string, resources []provider.Resource) *Provider {
	return &Provider{
		name:      name,
		resources: resources,
	}
}

// NewFailingProvider creates a mock provider that always fails on Fetch.
func NewFailingProvider(name string) *Provider {
	return &Provider{
		name:      name,
		failFetch: true,
	}
}

// Name returns the name of the mock provider.
func (m *Provider) Name() string {
	return m.name
}

// Fetch returns the preconfigured resources or an error if configured to fail.
func (m *Provider) Fetch() ([]provider.Resource, error) {
	if m.failFetch {
		return nil, fmt.Errorf("mock provider %q: simulated fetch failure", m.name)
	}
	return m.resources, nil
}

// SetResources replaces the provider's resource list.
func (m *Provider) SetResources(resources []provider.Resource) {
	m.resources = resources
}
