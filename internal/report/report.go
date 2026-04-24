package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/user/driftwatch/internal/drift"
)

// Format represents the output format for a drift report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Writer renders drift results to an output stream.
type Writer struct {
	format Format
	out    io.Writer
}

// New creates a new report Writer.
func New(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write renders the drift results according to the configured format.
func (w *Writer) Write(results []drift.Result) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(results)
	default:
		return w.writeText(results)
	}
}

func (w *Writer) writeJSON(results []drift.Result) error {
	enc := json.NewEncoder(w.out)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func (w *Writer) writeText(results []drift.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(w.out, "✓ No drift detected.")
		return nil
	}

	tw := tabwriter.NewWriter(w.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE\tSTATUS\tDETAIL")
	fmt.Fprintln(tw, strings.Repeat("-", 60))

	for _, r := range results {
		detail := r.Detail
		if detail == "" {
			detail = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.ResourceID, r.Status, detail)
	}

	return tw.Flush()
}
