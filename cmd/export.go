package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/export"
	"github.com/driftwatch/internal/snapshot"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var (
	exportFormat       string
	exportStatePath    string
	exportSnapshotPath string
	exportOutput       string
)

func init() {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export drift results to a file",
		Long:  "Compare declared state against a snapshot and export drift results as CSV or JSON.",
		RunE:  runExport,
	}

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "Output format: csv or json")
	exportCmd.Flags().StringVar(&exportStatePath, "state", "", "Path to declared state file")
	exportCmd.Flags().StringVar(&exportSnapshotPath, "snapshot", "", "Path to snapshot file")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (defaults to stdout)")

	_ = exportCmd.MarkFlagRequired("state")
	_ = exportCmd.MarkFlagRequired("snapshot")

	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	declared, err := state.Load(exportStatePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	snap, err := snapshot.Load(exportSnapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	results := drift.Detect(declared.Resources, snap.Resources)

	w := cmd.OutOrStdout()
	if exportOutput != "" {
		f, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	opts := export.Options{
		Format: export.Format(exportFormat),
		Writer: w,
	}

	if err := export.Write(results, opts); err != nil {
		return fmt.Errorf("exporting results: %w", err)
	}

	if exportOutput != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Results exported to %s\n", exportOutput)
	}
	return nil
}
