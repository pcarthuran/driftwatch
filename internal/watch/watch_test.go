package watch_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
	"github.com/driftwatch/internal/watch"
)

type mockFetcher struct {
	resources []state.Resource
}

func (m *mockFetcher) Fetch(_ context.Context) ([]state.Resource, error) {
	return m.resources, nil
}

func writeTempState(t *testing.T, resources []state.Resource) string {
	t.Helper()
	data, err := json.Marshal(map[string]interface{}{"resources": resources})
	if err != nil {
		t.Fatalf("marshal state: %v", err)
	}
	p := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatalf("write state: %v", err)
	}
	return p
}

func sampleResources() []state.Resource {
	return []state.Resource{
		{ID: "res-1", Type: "instance", Provider: "aws", Fields: map[string]interface{}{"size": "t2.micro"}},
	}
}

func TestRun_NoDrift(t *testing.T) {
	res := sampleResources()
	path := writeTempState(t, res)
	fetcher := &mockFetcher{resources: res}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := watch.Options{StatePath: path, Interval: 50 * time.Millisecond, MaxRuns: 2}
	ch := watch.Run(ctx, fetcher, opts)

	var results []watch.Result
	for r := range ch {
		results = append(results, r)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 runs, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error: %v", r.Err)
		}
		for _, dr := range r.Results {
			if dr.Status != drift.StatusClean {
				t.Errorf("expected clean, got %s", dr.Status)
			}
		}
	}
}

func TestRun_CancelStops(t *testing.T) {
	res := sampleResources()
	path := writeTempState(t, res)
	fetcher := &mockFetcher{resources: res}

	ctx, cancel := context.WithCancel(context.Background())
	opts := watch.Options{StatePath: path, Interval: 500 * time.Millisecond, MaxRuns: 0}
	ch := watch.Run(ctx, fetcher, opts)

	<-ch // first result
	cancel()

	// drain; channel must close
	for range ch {
	}
}
