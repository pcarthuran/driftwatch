package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/compare"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare <baseline> <current>",
	Short: "Compare two snapshot files and report drift between them",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

var compareOutput string

func init() {
	compareCmd.Flags().StringVarP(&compareOutput, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	baselinePath := args[0]
	currentPath := args[1]

	result, err := compare.Compare(baselinePath, currentPath)
	if err != nil {
		return fmt.Errorf("compare failed: %w", err)
	}

	switch compareOutput {
	case "json":
		if err := writeCompareJSON(result); err != nil {
			return err
		}
	default:
		compare.Write(os.Stdout, result)
	}

	if len(result.Diffs) > 0 {
		os.Exit(1)
	}
	return nil
}

func writeCompareJSON(r *compare.Result) error {
	enc := jsonEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("encoding JSON output: %w", err)
	}
	return nil
}
