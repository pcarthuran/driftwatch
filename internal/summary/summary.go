package summary

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Stats holds aggregated counts from a set of drift results.
type Stats struct {
	Total    int
	Missing  int
	Extra    int
	Modified int
	Clean    int
}

// Compute derives Stats from a slice of DetectResult.
func Compute(results []drift.DetectResult) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case drift.StatusMissing:
			s.Missing++
		case drift.StatusExtra:
			s.Extra++
		case drift.StatusModified:
			s.Modified++
		case drift.StatusOK:
			s.Clean++
		}
	}
	return s
}

// Write renders a human-readable summary table to w.
func Write(w io.Writer, results []drift.DetectResult) error {
	s := Compute(results)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SUMMARY")
	fmt.Fprintln(tw, "-------")
	fmt.Fprintf(tw, "Total resources:\t%d\n", s.Total)
	fmt.Fprintf(tw, "Clean:\t%d\n", s.Clean)
	fmt.Fprintf(tw, "Modified:\t%d\n", s.Modified)
	fmt.Fprintf(tw, "Missing:\t%d\n", s.Missing)
	fmt.Fprintf(tw, "Extra:\t%d\n", s.Extra)

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("summary: flush tabwriter: %w", err)
	}

	if s.Modified+s.Missing+s.Extra > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "DRIFTED RESOURCES")
		fmt.Fprintln(w, "-----------------")

		sorted := make([]drift.DetectResult, len(results))
		copy(sorted, results)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ResourceID < sorted[j].ResourceID
		})

		tw2 := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw2, "ID\tPROVIDER\tTYPE\tSTATUS")
		for _, r := range sorted {
			if r.Status != drift.StatusOK {
				fmt.Fprintf(tw2, "%s\t%s\t%s\t%s\n",
					r.ResourceID, r.Provider, r.ResourceType, r.Status)
			}
		}
		if err := tw2.Flush(); err != nil {
			return fmt.Errorf("summary: flush drifted table: %w", err)
		}
	}

	return nil
}
