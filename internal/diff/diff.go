package diff

import (
	"fmt"
	"sort"
	"strings"
)

// FieldDiff represents a single field-level difference between two values.
type FieldDiff struct {
	Field    string
	Expected interface{}
	Actual   interface{}
}

// String returns a human-readable representation of the field diff.
func (f FieldDiff) String() string {
	return fmt.Sprintf("  field %q: expected=%v actual=%v", f.Field, f.Expected, f.Actual)
}

// ResourceDiff holds the diff result for a single resource.
type ResourceDiff struct {
	ResourceID string
	Kind       string // "modified", "missing", "extra"
	Fields     []FieldDiff
}

// Summary returns a brief human-readable summary of the resource diff.
func (r ResourceDiff) Summary() string {
	switch r.Kind {
	case "missing":
		return fmt.Sprintf("[MISSING] %s", r.ResourceID)
	case "extra":
		return fmt.Sprintf("[EXTRA]   %s", r.ResourceID)
	case "modified":
		lines := []string{fmt.Sprintf("[MODIFIED] %s", r.ResourceID)}
		for _, f := range r.Fields {
			lines = append(lines, f.String())
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprintf("[UNKNOWN] %s", r.ResourceID)
}

// CompareFields computes field-level diffs between two generic maps.
func CompareFields(expected, actual map[string]interface{}) []FieldDiff {
	var diffs []FieldDiff
	keys := mergedKeys(expected, actual)
	for _, k := range keys {
		ev, eok := expected[k]
		av, aok := actual[k]
		if !eok {
			diffs = append(diffs, FieldDiff{Field: k, Expected: nil, Actual: av})
			continue
		}
		if !aok {
			diffs = append(diffs, FieldDiff{Field: k, Expected: ev, Actual: nil})
			continue
		}
		if fmt.Sprintf("%v", ev) != fmt.Sprintf("%v", av) {
			diffs = append(diffs, FieldDiff{Field: k, Expected: ev, Actual: av})
		}
	}
	return diffs
}

func mergedKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
