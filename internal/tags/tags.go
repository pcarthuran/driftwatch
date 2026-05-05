// Package tags provides utilities for filtering and matching resources
// based on tag key-value pairs declared in state or live snapshots.
package tags

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Filter holds tag-based filtering criteria.
type Filter struct {
	Required map[string]string // key=value pairs that must all match
	Keys     []string          // keys that must be present (any value)
}

// NewFilter parses a slice of raw tag expressions into a Filter.
// Supported forms: "key=value" and "key".
func NewFilter(exprs []string) (*Filter, error) {
	f := &Filter{
		Required: make(map[string]string),
	}
	for _, expr := range exprs {
		if expr == "" {
			continue
		}
		if idx := strings.IndexByte(expr, '='); idx >= 0 {
			key := strings.TrimSpace(expr[:idx])
			val := strings.TrimSpace(expr[idx+1:])
			if key == "" {
				return nil, fmt.Errorf("tags: empty key in expression %q", expr)
			}
			f.Required[key] = val
		} else {
			key := strings.TrimSpace(expr)
			if key == "" {
				return nil, fmt.Errorf("tags: empty key in expression %q", expr)
			}
			f.Keys = append(f.Keys, key)
		}
	}
	return f, nil
}

// Apply returns only the drift results whose resource ID appears in
// resources that satisfy the tag filter. tagMap maps resource ID to its tags.
func Apply(results []drift.ResourceDiff, tagMap map[string]map[string]string, f *Filter) []drift.ResourceDiff {
	if f == nil || (len(f.Required) == 0 && len(f.Keys) == 0) {
		return results
	}
	var out []drift.ResourceDiff
	for _, r := range results {
		tags := tagMap[r.ResourceID]
		if matches(tags, f) {
			out = append(out, r)
		}
	}
	return out
}

// matches returns true when the given tags satisfy all criteria in f.
func matches(tags map[string]string, f *Filter) bool {
	for k, v := range f.Required {
		got, ok := tags[k]
		if !ok || got != v {
			return false
		}
	}
	for _, k := range f.Keys {
		if _, ok := tags[k]; !ok {
			return false
		}
	}
	return true
}
