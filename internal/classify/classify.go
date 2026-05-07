// Package classify categorises drift results by severity level
// based on the type and number of field changes detected.
package classify

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Severity represents the urgency level of a drift finding.
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityLow      Severity = "LOW"
	SeverityNone     Severity = "NONE"
)

// Classification holds a drift result alongside its computed severity.
type Classification struct {
	Result   drift.ResourceResult
	Severity Severity
}

// Classify assigns a Severity to each ResourceResult and returns the list.
func Classify(results []drift.ResourceResult) []Classification {
	out := make([]Classification, 0, len(results))
	for _, r := range results {
		out = append(out, Classification{
			Result:   r,
			Severity: severityFor(r),
		})
	}
	return out
}

// severityFor derives a Severity from a single ResourceResult.
func severityFor(r drift.ResourceResult) Severity {
	switch {
	case r.Status == drift.StatusMissing:
		return SeverityCritical
	case r.Status == drift.StatusExtra:
		return SeverityHigh
	case r.Status == drift.StatusModified && len(r.Diffs) >= 3:
		return SeverityHigh
	case r.Status == drift.StatusModified:
		return SeverityLow
	default:
		return SeverityNone
	}
}

// Write renders a classified results table to w.
func Write(w io.Writer, classifications []Classification) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SEVERITY\tID\tTYPE\tPROVIDER\tSTATUS")
	fmt.Fprintln(tw, "--------\t--\t----\t--------\t------")
	for _, c := range classifications {
		r := c.Result
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			c.Severity,
			r.ResourceID,
			r.ResourceType,
			r.Provider,
			r.Status,
		)
	}
	return tw.Flush()
}
