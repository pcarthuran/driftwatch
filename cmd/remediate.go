package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/remediate"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

// remediateCmd defines the CLI surface for the remediate subcommand.
// It computes a remediation plan from the latest drift results and
// writes the suggested corrective actions to stdout or a file.
var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Generate a remediation plan for detected drift",
	Long: `Compares live infrastructure state against a declared state file
and produces a human-readable remediation plan that describes the
steps required to bring live resources back into compliance.

The plan is written to stdout by default. Use --output to write it
to a file instead.`,
	RunE: runRemediate,
}

func init() {
	remediateCmd.Flags().StringP("state", "s", "", "Path to the declared state file (JSON or YAML)")
	remediateCmd.Flags().StringP("snapshot", "n", "", "Path to a previously saved live-state snapshot (JSON or YAML)")
	remediateCmd.Flags().StringP("output", "o", "", "Write the remediation plan to this file instead of stdout")
	remediateCmd.Flags().BoolP("dry-run", "d", false, "Print the plan without writing to disk")

	_ = remediateCmd.MarkFlagRequired("state")
	_ = remediateCmd.MarkFlagRequired("snapshot")

	rootCmd.AddCommand(remediateCmd)
}

func runRemediate(cmd *cobra.Command, _ []string) error {
	statePath, _ := cmd.Flags().GetString("state")
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	outputPath, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Load declared state.
	declared, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("loading state file: %w", err)
	}

	// Load live snapshot.
	live, err := state.Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot file: %w", err)
	}

	// Run drift detection.
	results := drift.Detect(declared.Resources, live.Resources)

	// Build the remediation plan.
	plan := remediate.Plan(results)

	if len(plan) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No remediation required — infrastructure matches declared state.")
		return nil
	}

	// Determine output destination.
	if dryRun || outputPath == "" {
		// Write plan to stdout.
		if err := remediate.Write(plan, cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("writing remediation plan: %w", err)
		}
		return nil
	}

	// Write plan to the specified file.
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file %q: %w", outputPath, err)
	}
	defer f.Close()

	if err := remediate.Write(plan, f); err != nil {
		return fmt.Errorf("writing remediation plan to %q: %w", outputPath, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Remediation plan written to %s\n", outputPath)
	return nil
}
