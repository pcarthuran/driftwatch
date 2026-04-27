package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// Options configures the export behaviour.
type Options struct {
	Format Format
	Writer io.Writer
}

// Write serialises drift results to the configured format.
func Write(results []drift.Result, opts Options) error {
	switch opts.Format {
	case FormatCSV:
		return writeCSV(results, opts.Writer)
	case FormatJSON:
		return writeJSON(results, opts.Writer)
	default:
		return fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func writeCSV(results []drift.Result, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"resource_id", "provider", "type", "status", "field", "declared", "live"}); err != nil {
		return err
	}

	for _, r := range results {
		if len(r.Diffs) == 0 {
			if err := cw.Write([]string{r.ResourceID, r.Provider, r.Type, string(r.Status), "", "", ""}); err != nil {
				return err
			}
			continue
		}
		sortedDiffs := r.Diffs
		sort.Slice(sortedDiffs, func(i, j int) bool {
			return sortedDiffs[i].Field < sortedDiffs[j].Field
		})
		for _, d := range sortedDiffs {
			row := []string{
				r.ResourceID, r.Provider, r.Type, string(r.Status),
				d.Field, fmt.Sprintf("%v", d.Declared), fmt.Sprintf("%v", d.Live),
			}
			if err := cw.Write(row); err != nil {
				return err
			}
		}
	}

	cw.Flush()
	return cw.Error()
}

func writeJSON(results []drift.Result, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
