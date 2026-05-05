package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/score"
	"github.com/driftwatch/internal/snapshot"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var (
	scoreStateFile    string
	scoreSnapshotFile string
)

func init() {
	scoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Compute a drift health score for the current infrastructure state",
		Long: `Compares a snapshot of live resources against a declared state file
and outputs a health score from 0 to 100 with a letter grade.`,
		RunE: runScore,
	}

	scoreCmd.Flags().StringVarP(&scoreStateFile, "state", "s", "", "Path to declared state file (required)")
	scoreCmd.Flags().StringVarP(&scoreSnapshotFile, "snapshot", "n", "", "Path to live snapshot file (required)")
	_ = scoreCmd.MarkFlagRequired("state")
	_ = scoreCmd.MarkFlagRequired("snapshot")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, _ []string) error {
	declared, err := state.Load(scoreStateFile)
	if err != nil {
		return fmt.Errorf("loading state file: %w", err)
	}

	live, err := snapshot.Load(scoreSnapshotFile)
	if err != nil {
		return fmt.Errorf("loading snapshot file: %w", err)
	}

	results, err := drift.Detect(declared.Resources, live.Resources)
	if err != nil {
		return fmt.Errorf("detecting drift: %w", err)
	}

	r := score.Compute(results)
	if err := score.Write(os.Stdout, r); err != nil {
		return fmt.Errorf("writing score: %w", err)
	}

	if r.Drifted > 0 {
		os.Exit(1)
	}
	return nil
}
