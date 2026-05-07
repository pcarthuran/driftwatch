package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/group"
	"github.com/driftwatch/internal/snapshot"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var (
	groupBy       string
	groupSnapshot string
	groupState    string
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Group drift results by provider, type, or status",
		RunE:  runGroup,
	}
	groupCmd.Flags().StringVar(&groupBy, "by", "provider", "Grouping dimension: provider | type | status")
	groupCmd.Flags().StringVar(&groupSnapshot, "snapshot", "", "Path to a saved snapshot file")
	groupCmd.Flags().StringVar(&groupState, "state", "", "Path to a declared state file")
	_ = groupCmd.MarkFlagRequired("snapshot")
	_ = groupCmd.MarkFlagRequired("state")
	rootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, args []string) error {
	snap, err := snapshot.Load(groupSnapshot)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	declared, err := state.Load(groupState)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	results := drift.Detect(snap.Resources, declared.Resources)

	var dim group.By
	switch groupBy {
	case "type":
		dim = group.ByType
	case "status":
		dim = group.ByStatus
	default:
		dim = group.ByProvider
	}

	groups := group.Compute(results, dim)
	group.Write(os.Stdout, groups)
	return nil
}
