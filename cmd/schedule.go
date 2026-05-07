package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/driftwatch/internal/schedule"
	"github.com/spf13/cobra"
)

func init() {
	var (
		scheduleID       string
		scheduleName     string
		scheduleProvider string
		scheduleState    string
		scheduleInterval time.Duration
	)

	addCmd := &cobra.Command{
		Use:   "schedule-add",
		Short: "Add a drift-check schedule entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := schedule.NewStore(".driftwatch/schedules")
			if err != nil {
				return err
			}
			e := schedule.Entry{
				ID:        scheduleID,
				Name:      scheduleName,
				Provider:  scheduleProvider,
				StateFile: scheduleState,
				Interval:  scheduleInterval,
				Enabled:   true,
			}
			if err := store.Save(e); err != nil {
				return err
			}
			fmt.Printf("Schedule %q saved.\n", scheduleID)
			return nil
		},
	}
	addCmd.Flags().StringVar(&scheduleID, "id", "", "Unique schedule ID (required)")
	addCmd.Flags().StringVar(&scheduleName, "name", "", "Human-readable name")
	addCmd.Flags().StringVar(&scheduleProvider, "provider", "", "Provider name (aws|gcp|azure)")
	addCmd.Flags().StringVar(&scheduleState, "state", "", "Path to state file")
	addCmd.Flags().DurationVar(&scheduleInterval, "interval", 24*time.Hour, "Check interval (e.g. 6h, 24h)")
	_ = addCmd.MarkFlagRequired("id")

	listCmd := &cobra.Command{
		Use:   "schedule-list",
		Short: "List all drift-check schedule entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := schedule.NewStore(".driftwatch/schedules")
			if err != nil {
				return err
			}
			entries, err := store.List()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				fmt.Println("No schedules defined.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tPROVIDER\tINTERVAL\tENABLED")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n",
					e.ID, e.Name, e.Provider, e.Interval, e.Enabled)
			}
			return w.Flush()
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "schedule-delete <id>",
		Short: "Delete a drift-check schedule entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := schedule.NewStore(".driftwatch/schedules")
			if err != nil {
				return err
			}
			if err := store.Delete(args[0]); err != nil {
				return fmt.Errorf("deleting schedule %q: %w", args[0], err)
			}
			fmt.Printf("Schedule %q deleted.\n", args[0])
			return nil
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, deleteCmd)
}
