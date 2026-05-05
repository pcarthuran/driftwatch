// Package score computes a drift health score for a set of detection results.
package score

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Result holds the computed health score and its breakdown.
type Result struct {
	Total     int
	Clean     int
	Drifted   int
	Missing   int
	Extra     int
	Modified  int
	Score     float64 // 0.0 – 100.0
	Grade     string
}

// Compute derives a health score from a slice of drift detection results.
// Score = (clean / total) * 100, rounded to two decimal places.
// Returns a zero-value Result with Score 100 when results is empty.
func Compute(results []drift.ResourceResult) Result {
	if len(results) == 0 {
		return Result{Score: 100.0, Grade: "A"}
	}

	r := Result{Total: len(results)}
	for _, res := range results {
		switch res.Status {
		case drift.StatusMatch:
			r.Clean++
		case drift.StatusMissing:
			r.Missing++
			r.Drifted++
		case drift.StatusExtra:
			r.Extra++
			r.Drifted++
		case drift.StatusModified:
			r.Modified++
			r.Drifted++
		}
	}

	r.Score = float64(r.Clean) / float64(r.Total) * 100.0
	r.Grade = grade(r.Score)
	return r
}

// Write renders the score result as a human-readable table to w.
func Write(w io.Writer, r Result) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "=== Drift Health Score ===")
	fmt.Fprintf(tw, "Grade:\t%s\n", r.Grade)
	fmt.Fprintf(tw, "Score:\t%.2f / 100\n", r.Score)
	fmt.Fprintln(tw, "---")
	fmt.Fprintf(tw, "Total resources:\t%d\n", r.Total)
	fmt.Fprintf(tw, "Clean:\t%d\n", r.Clean)
	fmt.Fprintf(tw, "Drifted:\t%d\n", r.Drifted)
	fmt.Fprintf(tw, "  Missing:\t%d\n", r.Missing)
	fmt.Fprintf(tw, "  Extra:\t%d\n", r.Extra)
	fmt.Fprintf(tw, "  Modified:\t%d\n", r.Modified)
	return tw.Flush()
}

func grade(score float64) string {
	switch {
	case score >= 95:
		return "A"
	case score >= 80:
		return "B"
	case score >= 65:
		return "C"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}
