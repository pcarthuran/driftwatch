// Package ignore provides functionality for loading and evaluating
// drift ignore rules, allowing users to suppress known or accepted drift.
package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Rule represents a single ignore rule that can match resources by provider,
// type, and/or ID using glob patterns.
type Rule struct {
	Provider string `yaml:"provider" json:"provider"`
	Type     string `yaml:"type"     json:"type"`
	ID       string `yaml:"id"       json:"id"`
}

// RuleSet holds a collection of ignore rules.
type RuleSet struct {
	Rules []Rule
}

// LoadFile reads ignore rules from a .driftignore file.
// Each non-blank, non-comment line is parsed as: provider/type/id
// Any segment may be "*" or a glob pattern.
func LoadFile(path string) (*RuleSet, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &RuleSet{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var rules []Rule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rule, err := parseLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		rules = append(rules, rule)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &RuleSet{Rules: rules}, nil
}

// Matches reports whether the given provider, resourceType, and id are
// suppressed by any rule in the set.
func (rs *RuleSet) Matches(provider, resourceType, id string) bool {
	for _, r := range rs.Rules {
		if globMatch(r.Provider, provider) &&
			globMatch(r.Type, resourceType) &&
			globMatch(r.ID, id) {
			return true
		}
	}
	return false
}

// DefaultPath returns the conventional .driftignore path relative to dir.
func DefaultPath(dir string) string {
	return filepath.Join(dir, ".driftignore")
}

func parseLine(line string) (Rule, error) {
	parts := strings.SplitN(line, "/", 3)
	r := Rule{Provider: "*", Type: "*", ID: "*"}
	if len(parts) > 0 && parts[0] != "" {
		r.Provider = parts[0]
	}
	if len(parts) > 1 && parts[1] != "" {
		r.Type = parts[1]
	}
	if len(parts) > 2 && parts[2] != "" {
		r.ID = parts[2]
	}
	return r, nil
}

func globMatch(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	matched, err := doublestar.Match(pattern, value)
	if err != nil {
		return false
	}
	return matched
}
