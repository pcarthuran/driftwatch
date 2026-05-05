package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/watch"
	"github.com/spf13/cobra"
)

var (
	watchInterval string
	watchMaxRuns  int
	watchProvider string
	watchState    string
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Continuously watch for configuration drift",
		Long:  "Poll live infrastructure at a fixed interval and report drift against the declared state file.",
		RunE:  runWatch,
	}

	watchCmd.Flags().StringVarP(&watchState, "state", "s", "", "Path to state file (required)")
	watchCmd.Flags().StringVarP(&watchProvider, "provider", "p", "", "Provider name (required)")
	watchCmd.Flags().StringVarP(&watchInterval, "interval", "i", "60s", "Poll interval (e.g. 30s, 5m)")
	watchCmd.Flags().IntVar(&watchMaxRuns, "max-runs", 0, "Stop after N runs (0 = unlimited)")
	_ = watchCmd.MarkFlagRequired("state")
	_ = watchCmd.MarkFlagRequired("provider")

	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, _ []string) error {
	interval, err := time.ParseDuration(watchInterval)
	if err != nil {
		return fmt.Errorf("invalid interval %q: %w", watchInterval, err)
	}

	reg := provider.NewRegistry()
	p, err := reg.Get(watchProvider)
	if err != nil {
		return fmt.Errorf("provider %q not found: %w", watchProvider, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	opts := watch.Options{
		StatePath: watchState,
		Provider:  watchProvider,
		Interval:  interval,
		MaxRuns:   watchMaxRuns,
		Out:       cmd.OutOrStdout(),
	}

	for result := range watch.Run(ctx, p, opts) {
		if result.Err != nil {
			fmt.Fprintf(os.Stderr, "[watch] error: %v\n", result.Err)
			continue
		}
		driftCount := 0
		for _, r := range result.Results {
			if r.Status != drift.StatusClean {
				driftCount++
			}
		}
		fmt.Fprintf(cmd.OutOrStdout(), "[watch] run %d: %d drifted resource(s)\n", result.Run, driftCount)
	}
	return nil
}
