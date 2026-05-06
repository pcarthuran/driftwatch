// Package dedupe provides utilities for deduplicating drift detection results
// across multiple runs or providers, merging results by resource identity.
package dedupe

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Key uniquely identifies a resource across providers.
type Key struct {
	Provider string
	Type     string
	ID       string
}

// Deduplicate merges multiple slices of drift results, keeping the last-seen
// entry for each unique (provider, type, id) combination. Results are returned
// in a stable, sorted order.
func Deduplicate(sets ...[]drift.Result) []drift.Result {
	seen := make(map[Key]drift.Result)

	for _, results := range sets {
		for _, r := range results {
			k := Key{
				Provider: r.Provider,
				Type:     r.Type,
				ID:       r.ID,
			}
			seen[k] = r
		}
	}

	out := make([]drift.Result, 0, len(seen))
	for _, v := range seen {
		out = append(out, v)
	}

	sort.Slice(out, func(i, j int) bool {
		ki := fmt.Sprintf("%s/%s/%s", out[i].Provider, out[i].Type, out[i].ID)
		kj := fmt.Sprintf("%s/%s/%s", out[j].Provider, out[j].Type, out[j].ID)
		return ki < kj
	})

	return out
}

// Write renders a deduplicated result summary to w.
func Write(w io.Writer, results []drift.Result) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROVIDER\tTYPE\tID\tSTATUS")
	for _, r := range results {
		status := "clean"
		if r.Drifted {
			status = "drifted"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Provider, r.Type, r.ID, status)
	}
	tw.Flush()
}
