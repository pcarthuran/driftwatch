package trend

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/driftwatch/internal/baseline"
)

// Point represents drift counts at a point in time.
type Point struct {
	Timestamp time.Time
	Missing   int
	Extra     int
	Modified  int
	Total     int
}

// Report holds a time-ordered series of drift data points.
type Report struct {
	Points []Point
}

// Compute builds a trend Report from a slice of baseline entries.
func Compute(entries []baseline.Entry) Report {
	points := make([]Point, 0, len(entries))
	for _, e := range entries {
		missing, extra, modified := 0, 0, 0
		for _, r := range e.Results {
			switch r.Status {
			case "missing":
				missing++
			case "extra":
				extra++
			case "modified":
				modified++
			}
		}
		points = append(points, Point{
			Timestamp: e.SavedAt,
			Missing:   missing,
			Extra:     extra,
			Modified:  modified,
			Total:     missing + extra + modified,
		})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})
	return Report{Points: points}
}

// Write renders the trend report as a human-readable table to w.
func Write(w io.Writer, r Report) error {
	if len(r.Points) == 0 {
		_, err := fmt.Fprintln(w, "No trend data available.")
		return err
	}
	_, err := fmt.Fprintf(w, "%-26s %8s %8s %8s %8s\n", "Timestamp", "Missing", "Extra", "Modified", "Total")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "--------------------------------------------------------------------------")
	if err != nil {
		return err
	}
	for _, p := range r.Points {
		_, err = fmt.Fprintf(w, "%-26s %8d %8d %8d %8d\n",
			p.Timestamp.Format(time.RFC3339),
			p.Missing, p.Extra, p.Modified, p.Total)
		if err != nil {
			return err
		}
	}
	return nil
}
