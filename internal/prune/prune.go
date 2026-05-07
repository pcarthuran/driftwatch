// Package prune removes drift results that have been consistently clean
// across a configurable number of consecutive history entries.
package prune

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Result holds the outcome of a prune operation.
type Result struct {
	Retained []drift.ResourceResult
	Pruned   []drift.ResourceResult
}

// Prune filters out resources that have been clean (no drift) for at least
// minCleanRuns consecutive evaluations. Resources with drift are always retained.
func Prune(results []drift.ResourceResult, history [][]drift.ResourceResult, minCleanRuns int) Result {
	if minCleanRuns <= 0 {
		return Result{Retained: results}
	}

	// Build a map of resourceID -> consecutive clean run count from history.
	cleanCounts := make(map[string]int)
	for _, run := range history {
		for _, r := range run {
			if !r.Drifted {
				cleanCounts[r.ID]++
			} else {
				// Reset on any drift detection.
				cleanCounts[r.ID] = 0
			}
		}
	}

	var retained, pruned []drift.ResourceResult
	for _, r := range results {
		if r.Drifted {
			retained = append(retained, r)
			continue
		}
		if cleanCounts[r.ID] >= minCleanRuns {
			pruned = append(pruned, r)
		} else {
			retained = append(retained, r)
		}
	}

	sort.Slice(retained, func(i, j int) bool { return retained[i].ID < retained[j].ID })
	sort.Slice(pruned, func(i, j int) bool { return pruned[i].ID < pruned[j].ID })

	return Result{Retained: retained, Pruned: pruned}
}

// Write renders a prune result summary to w.
func Write(w io.Writer, r Result) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "=== Prune Summary ===")
	fmt.Fprintf(tw, "Retained:\t%d\n", len(r.Retained))
	fmt.Fprintf(tw, "Pruned:\t%d\n", len(r.Pruned))
	if len(r.Pruned) > 0 {
		fmt.Fprintln(tw, "\nPruned Resources:")
		fmt.Fprintln(tw, "  ID\tProvider\tType")
		for _, res := range r.Pruned {
			fmt.Fprintf(tw, "  %s\t%s\t%s\n", res.ID, res.Provider, res.Type)
		}
	}
	return tw.Flush()
}
