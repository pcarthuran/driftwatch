package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/graph"
	"github.com/driftwatch/internal/snapshot"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var (
	graphStateFile    string
	graphSnapshotFile string
	graphOutput       string
)

func init() {
	graphCmd := &cobra.Command{
		Use:   "graph",
		Short: "Render a DOT dependency graph of drifted resources",
		Long: `Compares a state file against a snapshot and writes a Graphviz DOT
graph to stdout (or --output file). Drifted nodes are highlighted in salmon;
clean nodes in lightblue.`,
		RunE: runGraph,
	}

	graphCmd.Flags().StringVarP(&graphStateFile, "state", "s", "", "path to declared state file (required)")
	graphCmd.Flags().StringVarP(&graphSnapshotFile, "snapshot", "n", "", "path to live snapshot file (required)")
	graphCmd.Flags().StringVarP(&graphOutput, "output", "o", "", "write DOT output to file instead of stdout")
	_ = graphCmd.MarkFlagRequired("state")
	_ = graphCmd.MarkFlagRequired("snapshot")

	rootCmd.AddCommand(graphCmd)
}

func runGraph(cmd *cobra.Command, args []string) error {
	declared, err := state.Load(graphStateFile)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	snap, err := snapshot.Load(graphSnapshotFile)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	results := drift.Detect(declared.Resources, snap.Resources)
	g := graph.Build(results)

	w := cmd.OutOrStdout()
	if graphOutput != "" {
		f, err := os.Create(graphOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return graph.Write(g, w)
}
