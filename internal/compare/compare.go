// Package compare provides utilities for comparing two snapshots of
// infrastructure resources and producing a structured diff summary.
package compare

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/snapshot"
)

// Result holds the outcome of comparing two named snapshots.
type Result struct {
	BaselineLabel string
	CurrentLabel  string
	Diffs         []drift.ResourceDiff
}

// Compare loads two snapshot files and returns a Result containing all
// detected drift between them.
func Compare(baselinePath, currentPath string) (*Result, error) {
	base, err := snapshot.Load(baselinePath)
	if err != nil {
		return nil, fmt.Errorf("loading baseline snapshot %q: %w", baselinePath, err)
	}

	current, err := snapshot.Load(currentPath)
	if err != nil {
		return nil, fmt.Errorf("loading current snapshot %q: %w", currentPath, err)
	}

	diffs := drift.Detect(base.Resources, current.Resources)

	return &Result{
		BaselineLabel: baselinePath,
		CurrentLabel:  currentPath,
		Diffs:         diffs,
	}, nil
}

// Write renders the comparison result as a human-readable table to w.
func Write(w io.Writer, r *Result) {
	fmt.Fprintf(w, "Comparing snapshots:\n")
	fmt.Fprintf(w, "  baseline : %s\n", r.BaselineLabel)
	fmt.Fprintf(w, "  current  : %s\n\n", r.CurrentLabel)

	if len(r.Diffs) == 0 {
		fmt.Fprintln(w, "No drift detected between snapshots.")
		return
	}

	fmt.Fprintf(w, "%d resource(s) differ:\n\n", len(r.Diffs))

	sort.Slice(r.Diffs, func(i, j int) bool {
		return r.Diffs[i].ResourceID < r.Diffs[j].ResourceID
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE\tSTATUS\tDETAIL")
	for _, d := range r.Diffs {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", d.ResourceID, d.Status, d.Summary())
	}
	tw.Flush()
}
