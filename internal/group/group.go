// Package group provides grouping of drift results by a chosen dimension.
package group

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// By controls the grouping dimension.
type By string

const (
	ByProvider By = "provider"
	ByType     By = "type"
	ByStatus   By = "status"
)

// Group holds a named bucket of drift results.
type Group struct {
	Key     string
	Results []drift.Result
}

// Compute partitions results into groups based on the chosen dimension.
// Unknown dimension values fall back to ByProvider.
func Compute(results []drift.Result, by By) []Group {
	buckets := make(map[string][]drift.Result)

	for _, r := range results {
		var key string
		switch by {
		case ByType:
			key = r.Resource.Type
		case ByStatus:
			if r.Drifted {
				key = "drifted"
			} else {
				key = "clean"
			}
		default:
			key = r.Resource.Provider
		}
		if key == "" {
			key = "(unknown)"
		}
		buckets[key] = append(buckets[key], r)
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	groups := make([]Group, 0, len(keys))
	for _, k := range keys {
		groups = append(groups, Group{Key: k, Results: buckets[k]})
	}
	return groups
}

// Write renders the grouped results as a human-readable table.
func Write(w io.Writer, groups []Group) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "GROUP\tTOTAL\tDRIFTED\tCLEAN")
	for _, g := range groups {
		total := len(g.Results)
		drifted := 0
		for _, r := range g.Results {
			if r.Drifted {
				drifted++
			}
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\n", g.Key, total, drifted, total-drifted)
	}
	tw.Flush()
}
