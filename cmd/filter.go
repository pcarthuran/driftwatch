package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/filter"
	"github.com/driftwatch/internal/snapshot"
	"github.com/spf13/cobra"
)

var (
	filterProviders []string
	filterTypes     []string
	filterIDs       []string
	filterLabelKey  string
	filterLabelVal  string
	filterSnapshot  string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter resources from a snapshot by provider, type, ID, or label",
	RunE:  runFilter,
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().StringSliceVar(&filterProviders, "provider", nil, "Filter by provider (e.g. aws,gcp)")
	filterCmd.Flags().StringSliceVar(&filterTypes, "type", nil, "Filter by resource type")
	filterCmd.Flags().StringSliceVar(&filterIDs, "id", nil, "Filter by resource ID")
	filterCmd.Flags().StringVar(&filterLabelKey, "label-key", "", "Filter by attribute key")
	filterCmd.Flags().StringVar(&filterLabelVal, "label-val", "", "Filter by attribute value (requires --label-key)")
	filterCmd.Flags().StringVar(&filterSnapshot, "snapshot", "", "Path to snapshot file to filter")
	_ = filterCmd.MarkFlagRequired("snapshot")
}

func runFilter(cmd *cobra.Command, args []string) error {
	snap, err := snapshot.Load(filterSnapshot)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	opts := filter.Options{
		Providers: filterProviders,
		Types:     filterTypes,
		IDs:       filterIDs,
		LabelKey:  filterLabelKey,
		LabelVal:  filterLabelVal,
	}

	matched := filter.Apply(snap.Resources, opts)

	if len(matched) == 0 {
		fmt.Fprintln(os.Stdout, "No resources matched the given filters.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Matched %d resource(s):\n", len(matched))
	for _, r := range matched {
		fmt.Fprintf(os.Stdout, "  [%s] %s/%s\n", r.Provider, r.Type, r.ID)
	}
	return nil
}
