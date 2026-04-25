package provider

import (
	"fmt"
	"sort"

	"github.com/driftwatch/internal/state"
)

// Provider defines the interface for fetching live infrastructure resources.
type Provider interface {
	Name() string
	FetchResources() ([]state.Resource, error)
}

// Registry holds registered providers by name.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(p Provider) error {
	name := p.Name()
	if name == "" {
		return fmt.Errorf("provider name must not be empty")
	}
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %q is already registered", name)
	}
	r.providers[name] = p
	return nil
}

// Get retrieves a provider by name.
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q not found", name)
	}
	return p, nil
}

// Names returns a sorted list of registered provider names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// FetchAll fetches resources from all registered providers and merges them.
func (r *Registry) FetchAll() ([]state.Resource, error) {
	var all []state.Resource
	for _, name := range r.Names() {
		p := r.providers[name]
		resources, err := p.FetchResources()
		if err != nil {
			return nil, fmt.Errorf("provider %q fetch error: %w", name, err)
		}
		all = append(all, resources...)
	}
	return all, nil
}
