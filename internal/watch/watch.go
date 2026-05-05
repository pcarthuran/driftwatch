package watch

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
)

// Options configures a Watch run.
type Options struct {
	StatePath string
	Provider  string
	Interval  time.Duration
	MaxRuns   int // 0 = unlimited
	Out       io.Writer
}

// Result holds the outcome of a single watch tick.
type Result struct {
	At      time.Time
	Run     int
	Results []drift.ResourceResult
	Err     error
}

// Run watches for drift on a fixed interval, sending results to the returned
// channel. The channel is closed when ctx is cancelled or MaxRuns is reached.
func Run(ctx context.Context, fetcher drift.Fetcher, opts Options) <-chan Result {
	out := make(chan Result)
	w := opts.Out
	if w == nil {
		w = os.Stdout
	}

	go func() {
		defer close(out)
		run := 0
		for {
			run++
			result := tick(fetcher, opts.StatePath)
			result.Run = run
			fmt.Fprintf(w, "[watch] run %d at %s\n", run, result.At.Format(time.RFC3339))
			select {
			case out <- result:
			case <-ctx.Done():
				return
			}
			if opts.MaxRuns > 0 && run >= opts.MaxRuns {
				return
			}
			select {
			case <-time.After(opts.Interval):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func tick(fetcher drift.Fetcher, statePath string) Result {
	r := Result{At: time.Now()}
	desired, err := state.Load(statePath)
	if err != nil {
		r.Err = fmt.Errorf("load state: %w", err)
		return r
	}
	live, err := fetcher.Fetch(context.Background())
	if err != nil {
		r.Err = fmt.Errorf("fetch live: %w", err)
		return r
	}
	r.Results = drift.Detect(live, desired)
	return r
}
