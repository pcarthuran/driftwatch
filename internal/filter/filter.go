package filter

import (
	"strings"

	"github.com/driftwatch/internal/state"
)

// Options holds filtering criteria for resources.
type Options struct {
	Providers []string
	Types     []string
	IDs       []string
	LabelKey  string
	LabelVal  string
}

// Apply returns only the resources from snapshot that match all non-empty
// criteria in opts.
func Apply(resources []state.Resource, opts Options) []state.Resource {
	var out []state.Resource
	for _, r := range resources {
		if !matchesProvider(r, opts.Providers) {
			continue
		}
		if !matchesType(r, opts.Types) {
			continue
		}
		if !matchesID(r, opts.IDs) {
			continue
		}
		if !matchesLabel(r, opts.LabelKey, opts.LabelVal) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func matchesProvider(r state.Resource, providers []string) bool {
	if len(providers) == 0 {
		return true
	}
	for _, p := range providers {
		if strings.EqualFold(r.Provider, p) {
			return true
		}
	}
	return false
}

func matchesType(r state.Resource, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if strings.EqualFold(r.Type, t) {
			return true
		}
	}
	return false
}

func matchesID(r state.Resource, ids []string) bool {
	if len(ids) == 0 {
		return true
	}
	for _, id := range ids {
		if r.ID == id {
			return true
		}
	}
	return false
}

func matchesLabel(r state.Resource, key, val string) bool {
	if key == "" {
		return true
	}
	v, ok := r.Attributes[key]
	if !ok {
		return false
	}
	if val == "" {
		return true
	}
	return fmt.Sprintf("%v", v) == val
}
