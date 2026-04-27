package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/diff"
	"github.com/driftwatch/internal/snapshot"
	"github.com/spf13/cobra"
)

var (
	diffBaseFile   string
	diffTargetFile string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show field-level differences between two snapshots",
	Long:  `Compare two snapshot files and display detailed field-level diffs for each resource.`,
	RunE:  runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&diffBaseFile, "base", "", "Path to the base (expected) snapshot file (required)")
	diffCmd.Flags().StringVar(&diffTargetFile, "target", "", "Path to the target (actual) snapshot file (required)")
	_ = diffCmd.MarkFlagRequired("base")
	_ = diffCmd.MarkFlagRequired("target")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	base, err := snapshot.Load(diffBaseFile)
	if err != nil {
		return fmt.Errorf("loading base snapshot: %w", err)
	}
	target, err := snapshot.Load(diffTargetFile)
	if err != nil {
		return fmt.Errorf("loading target snapshot: %w", err)
	}

	baseIndex := make(map[string]map[string]interface{})
	for _, r := range base.Resources {
		baseIndex[r.ID] = r.Attributes
	}
	targetIndex := make(map[string]map[string]interface{})
	for _, r := range target.Resources {
		targetIndex[r.ID] = r.Attributes
	}

	var diffs []diff.ResourceDiff

	for id, bAttrs := range baseIndex {
		if tAttrs, ok := targetIndex[id]; !ok {
			diffs = append(diffs, diff.ResourceDiff{ResourceID: id, Kind: "missing"})
		} else {
			fieldDiffs := diff.CompareFields(bAttrs, tAttrs)
			if len(fieldDiffs) > 0 {
				diffs = append(diffs, diff.ResourceDiff{ResourceID: id, Kind: "modified", Fields: fieldDiffs})
			}
		}
	}
	for id := range targetIndex {
		if _, ok := baseIndex[id]; !ok {
			diffs = append(diffs, diff.ResourceDiff{ResourceID: id, Kind: "extra"})
		}
	}

	if len(diffs) == 0 {
		fmt.Fprintln(os.Stdout, "No differences found between snapshots.")
		return nil
	}

	for _, d := range diffs {
		fmt.Fprintln(os.Stdout, d.Summary())
	}
	return nil
}
