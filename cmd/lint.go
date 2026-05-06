package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/lint"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var lintStatePath string

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Validate a state file for structural and semantic issues",
		RunE:  runLint,
	}
	lintCmd.Flags().StringVarP(&lintStatePath, "state", "s", "", "path to state file (required)")
	_ = lintCmd.MarkFlagRequired("state")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	snap, err := state.Load(lintStatePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	result := lint.Run(snap)
	lint.Write(os.Stdout, result)

	if result.HasErrors() {
		return fmt.Errorf("lint failed with %d error(s)", countErrors(result))
	}
	return nil
}

func countErrors(r *lint.Result) int {
	n := 0
	for _, f := range r.Findings {
		if f.Severity == lint.SeverityError {
			n++
		}
	}
	return n
}
