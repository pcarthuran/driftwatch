// Package flatten provides utilities for flattening nested drift results
// into a flat list of key-value field pairs suitable for tabular output.
package flatten

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Row represents a single flattened drift field entry.
type Row struct {
	Provider string
	ResourceID string
	ResourceType string
	Field string
	Wanted string
	Got string
	Status string
}

// Flatten converts a slice of drift.Result into a flat list of Rows,
// one row per diffed field. Resources with no drift are included with
// status "ok" and empty field columns.
func Flatten(results []drift.Result) []Row {
	var rows []Row
	for _, r := range results {
		if len(r.Diffs) == 0 {
			rows = append(rows, Row{
				Provider:     r.Provider,
				ResourceID:   r.ResourceID,
				ResourceType: r.ResourceType,
				Status:       "ok",
			})
			continue
		}
		for _, d := range r.Diffs {
			rows = append(rows, Row{
				Provider:     r.Provider,
				ResourceID:   r.ResourceID,
				ResourceType: r.ResourceType,
				Field:        d.Field,
				Wanted:       fmt.Sprintf("%v", d.Wanted),
				Got:          fmt.Sprintf("%v", d.Got),
				Status:       string(d.Kind),
			})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Provider != rows[j].Provider {
			return rows[i].Provider < rows[j].Provider
		}
		if rows[i].ResourceID != rows[j].ResourceID {
			return rows[i].ResourceID < rows[j].ResourceID
		}
		return rows[i].Field < rows[j].Field
	})
	return rows
}

// Write renders the flattened rows as a tab-separated table to w.
func Write(w io.Writer, rows []Row) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROVIDER\tID\tTYPE\tFIELD\tWANTED\tGOT\tSTATUS")
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.Provider, r.ResourceID, r.ResourceType,
			r.Field, r.Wanted, r.Got, r.Status)
	}
	return tw.Flush()
}
