package drift

import (
	"fmt"
	"sort"

	"github.com/driftwatch/internal/state"
)

// DriftResult represents a single detected drift between declared and live state.
type DriftResult struct {
	Resource string
	Field    string
	Declared interface{}
	Actual   interface{}
	Status   DriftStatus
}

// DriftStatus indicates the type of drift detected.
type DriftStatus string

const (
	StatusMissing  DriftStatus = "missing"   // resource exists in declared but not in live
	StatusExtra    DriftStatus = "extra"     // resource exists in live but not in declared
	StatusModified DriftStatus = "modified" // resource exists in both but values differ
)

// Report holds all drift results for a comparison run.
type Report struct {
	Results []DriftResult
	Drifted bool
}

// Detect compares a declared state snapshot against a live state snapshot
// and returns a Report describing any configuration drift.
func Detect(declared, live *state.Snapshot) (*Report, error) {
	if declared == nil {
		return nil, fmt.Errorf("declared snapshot must not be nil")
	}
	if live == nil {
		return nil, fmt.Errorf("live snapshot must not be nil")
	}

	report := &Report{}

	declaredMap := indexByID(declared.Resources)
	liveMap := indexByID(live.Resources)

	// Check for missing or modified resources.
	keys := sortedKeys(declaredMap)
	for _, id := range keys {
		dRes := declaredMap[id]
		lRes, exists := liveMap[id]
		if !exists {
			report.Results = append(report.Results, DriftResult{
				Resource: id,
				Status:   StatusMissing,
				Declared: dRes,
			})
			continue
		}
		if diffs := compareFields(id, dRes.Fields, lRes.Fields); len(diffs) > 0 {
			report.Results = append(report.Results, diffs...)
		}
	}

	// Check for extra resources in live state.
	for _, id := range sortedKeys(liveMap) {
		if _, exists := declaredMap[id]; !exists {
			report.Results = append(report.Results, DriftResult{
				Resource: id,
				Status:   StatusExtra,
				Actual:   liveMap[id],
			})
		}
	}

	report.Drifted = len(report.Results) > 0
	return report, nil
}

func indexByID(resources []state.Resource) map[string]state.Resource {
	m := make(map[string]state.Resource, len(resources))
	for _, r := range resources {
		m[r.ID] = r
	}
	return m
}

func sortedKeys(m map[string]state.Resource) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func compareFields(resourceID string, declared, actual map[string]interface{}) []DriftResult {
	var results []DriftResult
	for k, dv := range declared {
		av, exists := actual[k]
		if !exists || fmt.Sprintf("%v", dv) != fmt.Sprintf("%v", av) {
			results = append(results, DriftResult{
				Resource: resourceID,
				Field:    k,
				Declared: dv,
				Actual:   av,
				Status:   StatusModified,
			})
		}
	}
	return results
}
