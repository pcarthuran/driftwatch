package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule defines a single drift policy rule.
type Rule struct {
	ID       string            `json:"id" yaml:"id"`
	Provider string            `json:"provider" yaml:"provider"`
	Type     string            `json:"type" yaml:"type"`
	Field    string            `json:"field" yaml:"field"`
	Severity string            `json:"severity" yaml:"severity"` // info, warning, error
	Message  string            `json:"message" yaml:"message"`
	Labels   map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Policy holds a collection of rules.
type Policy struct {
	Rules []Rule `json:"rules" yaml:"rules"`
}

// Violation represents a rule that was triggered during evaluation.
type Violation struct {
	Rule       Rule
	ResourceID string
	Detail     string
}

// LoadFile reads a policy file (JSON or YAML) from disk.
func LoadFile(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read file: %w", err)
	}

	var p Policy
	switch {
	case strings.HasSuffix(path, ".json"):
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("policy: parse json: %w", err)
		}
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		if err := yaml.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("policy: parse yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("policy: unsupported format for %q", path)
	}

	if err := validate(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Evaluate checks drift results against the policy rules and returns violations.
func (p *Policy) Evaluate(driftedResources []DriftContext) []Violation {
	var violations []Violation
	for _, ctx := range driftedResources {
		for _, rule := range p.Rules {
			if !ruleMatches(rule, ctx) {
				continue
			}
			violations = append(violations, Violation{
				Rule:       rule,
				ResourceID: ctx.ResourceID,
				Detail:     fmt.Sprintf("field %q drifted on resource %s", rule.Field, ctx.ResourceID),
			})
		}
	}
	return violations
}

// DriftContext carries the minimal context needed for policy evaluation.
type DriftContext struct {
	ResourceID string
	Provider   string
	Type       string
	DriftedFields []string
	Labels     map[string]string
}

func ruleMatches(r Rule, ctx DriftContext) bool {
	if r.Provider != "" && r.Provider != ctx.Provider {
		return false
	}
	if r.Type != "" && r.Type != ctx.Type {
		return false
	}
	if r.Field != "" {
		found := false
		for _, f := range ctx.DriftedFields {
			if f == r.Field {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for k, v := range r.Labels {
		if ctx.Labels[k] != v {
			return false
		}
	}
	return true
}

func validate(p *Policy) error {
	valid := map[string]bool{"info": true, "warning": true, "error": true}
	for i, r := range p.Rules {
		if r.ID == "" {
			return fmt.Errorf("policy: rule[%d] missing id", i)
		}
		if r.Severity != "" && !valid[r.Severity] {
			return fmt.Errorf("policy: rule %q has invalid severity %q", r.ID, r.Severity)
		}
	}
	return nil
}
