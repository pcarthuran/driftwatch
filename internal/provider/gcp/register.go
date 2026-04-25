package gcp

import (
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// Register adds the GCP provider to the given registry using opts.
// Expected keys in opts: "project_id".
func Register(registry *provider.Registry, opts map[string]string) error {
	projectID, ok := opts["project_id"]
	if !ok || projectID == "" {
		return fmt.Errorf("gcp: missing required option 'project_id'")
	}

	p, err := New(projectID)
	if err != nil {
		return err
	}

	if err := registry.Register(p); err != nil {
		return fmt.Errorf("gcp: failed to register provider: %w", err)
	}

	return nil
}
