package aws

import (
	"context"
	"fmt"

	"github.com/user/driftwatch/internal/provider"
)

const ProviderName = "aws"

// Config holds AWS provider configuration.
type Config struct {
	Region    string            `yaml:"region" json:"region"`
	Profile   string            `yaml:"profile" json:"profile"`
	TagFilter map[string]string `yaml:"tag_filter" json:"tag_filter"`
}

// Provider implements provider.Provider for AWS resources.
type Provider struct {
	cfg     Config
	fetcher ResourceFetcher
}

// ResourceFetcher abstracts AWS API calls for testability.
type ResourceFetcher interface {
	FetchEC2Instances(ctx context.Context, region string) ([]provider.Resource, error)
}

// New creates a new AWS provider with the given config.
func New(cfg Config) *Provider {
	return &Provider{
		cfg:     cfg,
		fetcher: &defaultFetcher{region: cfg.Region, profile: cfg.Profile},
	}
}

// NewWithFetcher creates a new AWS provider with a custom fetcher (for testing).
func NewWithFetcher(cfg Config, fetcher ResourceFetcher) *Provider {
	return &Provider{cfg: cfg, fetcher: fetcher}
}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return ProviderName
}

// Fetch retrieves live AWS resources.
func (p *Provider) Fetch(ctx context.Context) ([]provider.Resource, error) {
	resources, err := p.fetcher.FetchEC2Instances(ctx, p.cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("aws provider fetch: %w", err)
	}
	return p.applyTagFilter(resources), nil
}

func (p *Provider) applyTagFilter(resources []provider.Resource) []provider.Resource {
	if len(p.cfg.TagFilter) == 0 {
		return resources
	}
	var filtered []provider.Resource
	for _, r := range resources {
		if matchesTags(r.Attributes, p.cfg.TagFilter) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func matchesTags(attrs map[string]interface{}, filter map[string]string) bool {
	for k, v := range filter {
		if val, ok := attrs[k]; !ok || fmt.Sprintf("%v", val) != v {
			return false
		}
	}
	return true
}
