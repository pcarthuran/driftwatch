package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
	"github.com/driftwatch/internal/snapshot"
	"github.com/driftwatch/internal/state"
	"github.com/driftwatch/internal/tags"
	"github.com/spf13/cobra"
)

var (
	tagsSnapshotFile string
	tagsStateFile    string
	tagsFormat       string
	tagsExprs        []string
)

func init() {
	tagsCmd := &cobra.Command{
		Use:   "tags",
		Short: "Filter drift results by resource tags",
		Long:  `Detect drift and show only resources matching the given tag expressions (key=value or key).`,
		RunE:  runTags,
	}

	tagsCmd.Flags().StringVarP(&tagsSnapshotFile, "snapshot", "s", "", "Path to snapshot file (required)")
	tagsCmd.Flags().StringVarP(&tagsStateFile, "state", "f", "", "Path to state file (required)")
	tagsCmd.Flags().StringVarP(&tagsFormat, "format", "o", "text", "Output format: text or json")
	tagsCmd.Flags().StringArrayVarP(&tagsExprs, "tag", "t", nil, "Tag expression: key=value or key (repeatable)")

	_ = tagsCmd.MarkFlagRequired("snapshot")
	_ = tagsCmd.MarkFlagRequired("state")

	rootCmd.AddCommand(tagsCmd)
}

func runTags(cmd *cobra.Command, _ []string) error {
	snap, err := snapshot.Load(tagsSnapshotFile)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	declared, err := state.Load(tagsStateFile)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	results := drift.Detect(snap.Resources, declared.Resources)

	if len(tagsExprs) > 0 {
		tagMap := buildTagMap(snap.Resources)
		f, err := tags.NewFilter(tagsExprs)
		if err != nil {
			return fmt.Errorf("parsing tag expressions: %w", err)
		}
		results = tags.Apply(results, tagMap, f)
	}

	rep := report.New(results)
	return rep.Write(os.Stdout, tagsFormat)
}

// buildTagMap extracts tags from each resource's Fields map.
func buildTagMap(resources []state.Resource) map[string]map[string]string {
	m := make(map[string]map[string]string, len(resources))
	for _, r := range resources {
		tagMap := make(map[string]string)
		if raw, ok := r.Fields["tags"]; ok {
			if kv, ok := raw.(map[string]interface{}); ok {
				for k, v := range kv {
					if s, ok := v.(string); ok {
						tagMap[k] = s
					}
				}
			}
		}
		m[r.ID] = tagMap
	}
	return m
}
