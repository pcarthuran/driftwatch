package policy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/policy"
)

func writeTempPolicy(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp policy: %v", err)
	}
	return path
}

func TestLoadFile_ValidYAML(t *testing.T) {
	content := `
rules:
  - id: rule-1
    provider: aws
    type: ec2
    field: instance_type
    severity: warning
    message: instance type changed
`
	path := writeTempPolicy(t, "policy.yaml", content)
	p, err := policy.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
	if p.Rules[0].ID != "rule-1" {
		t.Errorf("expected id rule-1, got %s", p.Rules[0].ID)
	}
}

func TestLoadFile_ValidJSON(t *testing.T) {
	content := `{"rules":[{"id":"r1","provider":"gcp","severity":"error"}]}`
	path := writeTempPolicy(t, "policy.json", content)
	p, err := policy.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Rules[0].Provider != "gcp" {
		t.Errorf("expected provider gcp")
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := policy.LoadFile("/nonexistent/policy.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_UnsupportedFormat(t *testing.T) {
	path := writeTempPolicy(t, "policy.toml", "")
	_, err := policy.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestLoadFile_InvalidSeverity(t *testing.T) {
	content := `rules:\n  - id: r1\n    severity: critical\n`
	path := writeTempPolicy(t, "policy.yaml", content)
	_, err := policy.LoadFile(path)
	if err == nil {
		t.Fatal("expected validation error for invalid severity")
	}
}

func sampleContexts() []policy.DriftContext {
	return []policy.DriftContext{
		{
			ResourceID:    "i-123",
			Provider:      "aws",
			Type:          "ec2",
			DriftedFields: []string{"instance_type", "ami"},
			Labels:        map[string]string{"env": "prod"},
		},
	}
}

func TestEvaluate_MatchingRule(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{ID: "r1", Provider: "aws", Type: "ec2", Field: "instance_type", Severity: "warning"},
		},
	}
	violations := p.Evaluate(sampleContexts())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].ResourceID != "i-123" {
		t.Errorf("unexpected resource id: %s", violations[0].ResourceID)
	}
}

func TestEvaluate_NoMatch(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{ID: "r1", Provider: "azure", Field: "sku", Severity: "info"},
		},
	}
	violations := p.Evaluate(sampleContexts())
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluate_LabelFilter(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{ID: "r1", Provider: "aws", Field: "ami", Severity: "error", Labels: map[string]string{"env": "prod"}},
		},
	}
	violations := p.Evaluate(sampleContexts())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for label match, got %d", len(violations))
	}
}
