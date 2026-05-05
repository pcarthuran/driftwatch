// Package watch provides continuous drift detection by polling live
// infrastructure at a configurable interval and comparing it against a
// declared state file.
//
// Basic usage:
//
//	opts := watch.Options{
//		StatePath: "infra/state.yaml",
//		Provider:  "aws",
//		Interval:  60 * time.Second,
//		MaxRuns:   0, // run forever
//	}
//	for result := range watch.Run(ctx, fetcher, opts) {
//		if result.Err != nil {
//			log.Printf("watch error: %v", result.Err)
//			continue
//		}
//		// process result.Results ...
//	}
//
// The channel returned by Run is closed when the context is cancelled or
// MaxRuns ticks have completed.
package watch
