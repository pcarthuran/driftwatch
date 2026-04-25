package azure

import (
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// Register adds the Azure provider to the given registry using the provided config.
// Expected keys in cfg: "subscription_id".
func Register(registry *provider.Registry, cfg map[string]string) error {
	subscriptionID, ok := cfg["subscription_id"]
	if !ok || subscriptionID == "" {
		return fmt.Errorf("azure: missing required config key: subscription_id")
	}

	p, err := New(subscriptionID)
	if err != nil {
		return err
	}

	if err := registry.Register(p); err != nil {
		return fmt.Errorf("azure: failed to register provider: %w", err)
	}

	return nil
}
