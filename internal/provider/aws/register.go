package aws

import (
	"fmt"

	"github.com/user/driftwatch/internal/provider"
)

// Register adds the AWS provider to the given registry using cfg.
// It returns an error if registration fails or the config is invalid.
func Register(registry *provider.Registry, cfg Config) error {
	if cfg.Region == "" {
		return fmt.Errorf("aws provider: region is required")
	}
	p := New(cfg)
	return registry.Register(p)
}
