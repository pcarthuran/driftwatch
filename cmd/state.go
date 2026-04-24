package cmd

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var stateFilePath string

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Inspect and validate a declared state file",
	Long:  `Load and display resources declared in a state file (.json or .yaml).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if stateFilePath == "" {
			return fmt.Errorf("--file flag is required")
		}

		sf, err := state.Load(stateFilePath)
		if err != nil {
			return fmt.Errorf("loading state: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "State version : %s\n", sf.Version)
		fmt.Fprintf(cmd.OutOrStdout(), "Resources     : %d\n", len(sf.Resources))
		fmt.Fprintln(cmd.OutOrStdout())

		for _, r := range sf.Resources {
			fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s (%s)\n", r.Provider, r.ID, r.Type)
			for k, v := range r.Attributes {
				fmt.Fprintf(cmd.OutOrStdout(), "      %s = %s\n", k, v)
			}
		}

		return nil
	},
}

func init() {
	stateCmd.Flags().StringVarP(&stateFilePath, "file", "f", "", "Path to the state file (.json or .yaml)")
	rootCmd.AddCommand(stateCmd)
}
