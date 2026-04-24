package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var (
	declaredPath string
	livePath     string
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect configuration drift between declared and live state",
	Long: `Compare a declared state file against a live state file and report
any configuration drift, including missing, extra, or modified resources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		declared, err := state.Load(declaredPath)
		if err != nil {
			return fmt.Errorf("loading declared state: %w", err)
		}

		live, err := state.Load(livePath)
		if err != nil {
			return fmt.Errorf("loading live state: %w", err)
		}

		report, err := drift.Detect(declared, live)
		if err != nil {
			return fmt.Errorf("running drift detection: %w", err)
		}

		if !report.Drifted {
			fmt.Println("✓ No drift detected.")
			return nil
		}

		fmt.Fprintf(os.Stderr, "✗ Drift detected (%d issue(s)):\n\n", len(report.Results))
		for _, r := range report.Results {
			switch r.Status {
			case drift.StatusMissing:
				fmt.Fprintf(os.Stderr, "  [MISSING]  %s\n", r.Resource)
			case drift.StatusExtra:
				fmt.Fprintf(os.Stderr, "  [EXTRA]    %s\n", r.Resource)
			case drift.StatusModified:
				fmt.Fprintf(os.Stderr, "  [MODIFIED] %s .%s: declared=%v actual=%v\n",
					r.Resource, r.Field, r.Declared, r.Actual)
			}
		}
		os.Exit(1)
		return nil
	},
}

func init() {
	detectCmd.Flags().StringVarP(&declaredPath, "declared", "d", "", "Path to declared state file (required)")
	detectCmd.Flags().StringVarP(&livePath, "live", "l", "", "Path to live state file (required)")
	_ = detectCmd.MarkFlagRequired("declared")
	_ = detectCmd.MarkFlagRequired("live")
	rootCmd.AddCommand(detectCmd)
}
