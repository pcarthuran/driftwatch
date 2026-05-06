// Package lint validates state file resources against structural and semantic
// rules before they are used in drift detection.
//
// Rules enforced:
//
//   - (error)   resource ID must not be empty
//   - (error)   resource type must not be empty
//   - (warning) provider must not be empty
//   - (warning) at least one field should be declared
//
// Usage:
//
//	snap, _ := state.Load("infra.yaml")
//	result := lint.Run(snap)
//	lint.Write(os.Stdout, result)
//	if result.HasErrors() {
//	    os.Exit(1)
//	}
package lint
