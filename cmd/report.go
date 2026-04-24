package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/report"
	"github.com/user/driftwatch/internal/state"
)

var (
	reportFormat string
	reportOutput string
)

func init() {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a drift report from a state file comparison",
		Long:  `Compares a declared state file against a live snapshot and outputs a formatted drift report.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runReport,
	}

	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "text", "Output format: text or json")
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "", "Write report to file instead of stdout")

	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	declaredPath := args[0]
	livePath := args[1]

	declared, err := state.Load(declaredPath)
	if err != nil {
		return fmt.Errorf("loading declared state: %w", err)
	}

	live, err := state.Load(livePath)
	if err != nil {
		return fmt.Errorf("loading live state: %w", err)
	}

	results := drift.Detect(declared, live)

	out := cmd.OutOrStdout()
	if reportOutput != "" {
		f, err := os.Create(reportOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	w := report.New(out, report.Format(reportFormat))
	return w.Write(results)
}
