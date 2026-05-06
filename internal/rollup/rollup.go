// Package rollup aggregates drift detection results across multiple providers
// into a single consolidated summary report.
package rollup

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// ProviderResult holds the drift detection results for a single provider.
type ProviderResult struct {
	Provider string
	Results  []drift.Result
}

// Report is the consolidated rollup across all providers.
type Report struct {
	Providers []ProviderResult
	Total     int
	Drifted   int
	Clean     int
	Missing   int
	Extra     int
	Modified  int
}

// Compute aggregates a slice of ProviderResults into a single Report.
func Compute(providers []ProviderResult) Report {
	r := Report{Providers: providers}
	for _, p := range providers {
		for _, res := range p.Results {
			r.Total++
			switch res.Status {
			case drift.StatusMatch:
				r.Clean++
			case drift.StatusMissing:
				r.Drifted++
				r.Missing++
			case drift.StatusExtra:
				r.Drifted++
				r.Extra++
			case drift.StatusModified:
				r.Drifted++
				r.Modified++
			}
		}
	}
	return r
}

// Write renders the rollup report as a human-readable table to w.
func Write(w io.Writer, r Report) error {
	fmt.Fprintln(w, "=== Drift Rollup Report ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROVIDER\tTOTAL\tCLEAN\tDRIFTED")

	// Sort providers for deterministic output.
	sorted := make([]ProviderResult, len(r.Providers))
	copy(sorted, r.Providers)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Provider < sorted[j].Provider
	})

	for _, p := range sorted {
		total, clean, drifted := 0, 0, 0
		for _, res := range p.Results {
			total++
			if res.Status == drift.StatusMatch {
				clean++
			} else {
				drifted++
			}
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\n", p.Provider, total, clean, drifted)
	}
	tw.Flush()

	fmt.Fprintf(w, "\nSummary: total=%d clean=%d drifted=%d (missing=%d extra=%d modified=%d)\n",
		r.Total, r.Clean, r.Drifted, r.Missing, r.Extra, r.Modified)
	return nil
}
