package remediate

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Action represents a suggested remediation action for a drifted resource.
type Action struct {
	ResourceID string
	Provider   string
	Kind       string // "missing", "extra", "modified"
	Suggestion string
}

// Plan generates a list of remediation actions from drift detection results.
func Plan(results []drift.Result) []Action {
	var actions []Action
	for _, r := range results {
		switch {
		case r.Missing:
			actions = append(actions, Action{
				ResourceID: r.ResourceID,
				Provider:   r.Provider,
				Kind:       "missing",
				Suggestion: fmt.Sprintf("Create resource '%s' (%s) to match declared state.", r.ResourceID, r.Provider),
			})
		case r.Extra:
			actions = append(actions, Action{
				ResourceID: r.ResourceID,
				Provider:   r.Provider,
				Kind:       "extra",
				Suggestion: fmt.Sprintf("Remove undeclared resource '%s' (%s) or add it to state.", r.ResourceID, r.Provider),
			})
		case len(r.Diffs) > 0:
			fields := make([]string, 0, len(r.Diffs))
			for _, d := range r.Diffs {
				fields = append(fields, fmt.Sprintf("%s (want %q, got %q)", d.Field, d.Expected, d.Actual))
			}
			actions = append(actions, Action{
				ResourceID: r.ResourceID,
				Provider:   r.Provider,
				Kind:       "modified",
				Suggestion: fmt.Sprintf("Update resource '%s' (%s): %s.", r.ResourceID, r.Provider, strings.Join(fields, "; ")),
			})
		}
	}
	return actions
}

// Write outputs the remediation plan in human-readable form to w.
func Write(w io.Writer, actions []Action) {
	if len(actions) == 0 {
		fmt.Fprintln(w, "No remediation actions required. Infrastructure matches declared state.")
		return
	}
	fmt.Fprintf(w, "Remediation Plan (%d action(s)):\n", len(actions))
	fmt.Fprintln(w, strings.Repeat("-", 50))
	for i, a := range actions {
		fmt.Fprintf(w, "%d. [%s] %s\n", i+1, strings.ToUpper(a.Kind), a.Suggestion)
	}
}
