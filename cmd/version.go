package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version   = "0.1.0"
	BuildDate = "unknown"
	Commit    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of driftwatch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("driftwatch v%s\n", Version)
		if verbose {
			fmt.Printf("  build date : %s\n", BuildDate)
			fmt.Printf("  commit     : %s\n", Commit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
