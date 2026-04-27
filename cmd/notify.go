package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/notify"
)

var (
	notifyChannel    string
	notifyWebhookURL string
	notifyTitle      string
)

func init() {
	notifyCmd := &cobra.Command{
		Use:   "notify",
		Short: "Send drift results to a notification channel",
		Long:  `Reads the latest baseline drift results and dispatches a summary to the configured webhook or Slack channel.`,
		RunE:  runNotify,
	}

	notifyCmd.Flags().StringVar(&notifyChannel, "channel", "webhook", "Notification channel: webhook or slack")
	notifyCmd.Flags().StringVar(&notifyWebhookURL, "webhook-url", "", "Destination webhook URL (required)")
	notifyCmd.Flags().StringVar(&notifyTitle, "title", "Drift Report", "Title for the notification message")
	notifyCmd.Flags().StringVar(&baselineDir, "baseline-dir", ".driftwatch/baselines", "Directory containing baseline results")

	_ = notifyCmd.MarkFlagRequired("webhook-url")

	rootCmd.AddCommand(notifyCmd)
}

func runNotify(cmd *cobra.Command, args []string) error {
	store, err := baseline.NewStore(baselineDir)
	if err != nil {
		return fmt.Errorf("notify: open baseline store: %w", err)
	}

	results, err := store.Latest()
	if err != nil {
		return fmt.Errorf("notify: load latest baseline: %w", err)
	}

	driftCount := 0
	for _, r := range results {
		if r.Status != "ok" {
			driftCount++
		}
	}

	msg := fmt.Sprintf("No drift detected across %d resources.", len(results))
	if driftCount > 0 {
		msg = fmt.Sprintf("%d of %d resources have drifted from declared state.", driftCount, len(results))
	}

	cfg := notify.Config{
		Channel:    notify.Channel(notifyChannel),
		WebhookURL: notifyWebhookURL,
	}
	sender := notify.New(cfg)

	payload := notify.Payload{
		Title:      notifyTitle,
		Message:    msg,
		DriftCount: driftCount,
	}

	if err := sender.Send(payload); err != nil {
		return fmt.Errorf("notify: send failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Notification sent (%s): %d drifted resources\n", notifyChannel, driftCount)
	return nil
}
