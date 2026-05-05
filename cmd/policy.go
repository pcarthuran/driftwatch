package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
	"github.com/driftwatch/internal/state"
	"github.com/spf13/cobra"
)

var (
	policyFile   string
	policyState  string
	policyFail   bool
)

func init() {
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Evaluate drift results against a policy file",
		Long:  `Load a policy file and evaluate detected drift against its rules, reporting violations.`,
		RunE:  runPolicy,
	}

	policyCmd.Flags().StringVarP(&policyFile, "policy", "p", "", "Path to policy file (JSON or YAML) (required)")
	policyCmd.Flags().StringVarP(&policyState, "state", "s", "", "Path to state file for drift detection (required)")
	policyCmd.Flags().BoolVar(&policyFail, "fail-on-violation", false, "Exit with non-zero code if violations are found")
	_ = policyCmd.MarkFlagRequired("policy")
	_ = policyCmd.MarkFlagRequired("state")

	rootCmd.AddCommand(policyCmd)
}

func runPolicy(cmd *cobra.Command, args []string) error {
	p, err := policy.LoadFile(policyFile)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}

	desired, err := state.Load(policyState)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	results, err := drift.Detect(nil, desired.Resources)
	if err != nil {
		return fmt.Errorf("detect drift: %w", err)
	}

	var contexts []policy.DriftContext
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		fields := make([]string, 0, len(r.Diffs))
		for _, d := range r.Diffs {
			fields = append(fields, d.Field)
		}
		contexts = append(contexts, policy.DriftContext{
			ResourceID:    r.ResourceID,
			Provider:      r.Provider,
			Type:          r.Type,
			DriftedFields: fields,
			Labels:        r.Labels,
		})
	}

	violations := p.Evaluate(contexts)

	if len(violations) == 0 {
		fmt.Println("No policy violations detected.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RULE ID\tSEVERITY\tRESOURCE\tDETAIL")
	for _, v := range violations {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.Rule.ID, v.Rule.Severity, v.ResourceID, v.Detail)
	}
	w.Flush()

	if policyFail && len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
