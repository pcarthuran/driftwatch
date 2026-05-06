// Package lint provides validation of state files against structural and
// semantic rules before drift detection is performed.
package lint

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/state"
)

// Severity indicates the level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Finding represents a single lint result for a resource.
type Finding struct {
	ResourceID string
	Field      string
	Message    string
	Severity   Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s.%s: %s", f.Severity, f.ResourceID, f.Field, f.Message)
}

// Result holds all findings from a lint run.
type Result struct {
	Findings []Finding
}

// HasErrors returns true if any finding has error severity.
func (r *Result) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Run validates all resources in the snapshot and returns a Result.
func Run(snap *state.Snapshot) *Result {
	result := &Result{}
	for _, res := range snap.Resources {
		if strings.TrimSpace(res.ID) == "" {
			result.Findings = append(result.Findings, Finding{
				ResourceID: "(unknown)",
				Field:      "id",
				Message:    "resource id must not be empty",
				Severity:   SeverityError,
			})
		}
		if strings.TrimSpace(res.Type) == "" {
			result.Findings = append(result.Findings, Finding{
				ResourceID: res.ID,
				Field:      "type",
				Message:    "resource type must not be empty",
				Severity:   SeverityError,
			})
		}
		if strings.TrimSpace(res.Provider) == "" {
			result.Findings = append(result.Findings, Finding{
				ResourceID: res.ID,
				Field:      "provider",
				Message:    "provider must not be empty",
				Severity:   SeverityWarning,
			})
		}
		if len(res.Fields) == 0 {
			result.Findings = append(result.Findings, Finding{
				ResourceID: res.ID,
				Field:      "fields",
				Message:    "resource has no fields defined",
				Severity:   SeverityWarning,
			})
		}
	}
	return result
}

// Write outputs findings to w in human-readable form.
func Write(w io.Writer, r *Result) {
	if len(r.Findings) == 0 {
		fmt.Fprintln(w, "lint: no issues found")
		return
	}
	fmt.Fprintf(w, "lint: %d issue(s) found\n", len(r.Findings))
	for _, f := range r.Findings {
		fmt.Fprintln(w, " ", f.String())
	}
}
