package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
	"github.com/driftwatch/internal/watch"
)

func TestIntegration_DriftDetectedAcrossRuns(t *testing.T) {
	desired := []state.Resource{
		{ID: "db-1", Type: "database", Provider: "gcp", Fields: map[string]interface{}{"tier": "db-f1-micro"}},
	}
	live := []state.Resource{
		{ID: "db-1", Type: "database", Provider: "gcp", Fields: map[string]interface{}{"tier": "db-n1-standard-1"}},
	}

	path := writeTempState(t, desired)
	fetcher := &mockFetcher{resources: live}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := watch.Options{StatePath: path, Interval: 20 * time.Millisecond, MaxRuns: 3}
	ch := watch.Run(ctx, fetcher, opts)

	var runs []watch.Result
	for r := range ch {
		runs = append(runs, r)
	}

	if len(runs) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(runs))
	}
	for i, r := range runs {
		if r.Err != nil {
			t.Errorf("run %d: unexpected error: %v", i+1, r.Err)
		}
		if len(r.Results) == 0 {
			t.Errorf("run %d: expected drift results", i+1)
			continue
		}
		if r.Results[0].Status != drift.StatusModified {
			t.Errorf("run %d: expected modified, got %s", i+1, r.Results[0].Status)
		}
	}
}

func TestIntegration_MaxRunsRespected(t *testing.T) {
	res := sampleResources()
	path := writeTempState(t, res)
	fetcher := &mockFetcher{resources: res}

	ctx := context.Background()
	opts := watch.Options{StatePath: path, Interval: 10 * time.Millisecond, MaxRuns: 5}
	ch := watch.Run(ctx, fetcher, opts)

	count := 0
	for range ch {
		count++
	}
	if count != 5 {
		t.Errorf("expected exactly 5 runs, got %d", count)
	}
}
