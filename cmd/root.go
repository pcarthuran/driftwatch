package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	verbose  bool
)

var rootCmd = &cobra.Command{
	Use:   "driftwatch",
	Short: "Detect configuration drift between live infrastructure and declared state",
	Long: `driftwatch compares your declared infrastructure state files against
live infrastructure to surface configuration drift early.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .driftwatch.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
