package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/ignore"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".driftignore")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp ignore file: %v", err)
	}
	return path
}

func TestLoadFile_MissingFile_ReturnsEmpty(t *testing.T) {
	rs, err := ignore.LoadFile("/nonexistent/.driftignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(rs.Rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rs.Rules))
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempIgnore(t, "# this is a comment\n\naws/ec2/i-12345\n")
	rs, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rs.Rules))
	}
}

func TestMatches_ExactRule(t *testing.T) {
	path := writeTempIgnore(t, "aws/ec2/i-12345\n")
	rs, _ := ignore.LoadFile(path)

	if !rs.Matches("aws", "ec2", "i-12345") {
		t.Error("expected match for exact rule")
	}
	if rs.Matches("aws", "ec2", "i-99999") {
		t.Error("expected no match for different id")
	}
}

func TestMatches_WildcardProvider(t *testing.T) {
	path := writeTempIgnore(t, "*/ec2/i-12345\n")
	rs, _ := ignore.LoadFile(path)

	if !rs.Matches("aws", "ec2", "i-12345") {
		t.Error("expected wildcard provider to match aws")
	}
	if !rs.Matches("gcp", "ec2", "i-12345") {
		t.Error("expected wildcard provider to match gcp")
	}
}

func TestMatches_WildcardAll(t *testing.T) {
	path := writeTempIgnore(t, "*\n")
	rs, _ := ignore.LoadFile(path)

	if !rs.Matches("aws", "s3", "my-bucket") {
		t.Error("expected catch-all rule to match any resource")
	}
}

func TestMatches_GlobPattern(t *testing.T) {
	path := writeTempIgnore(t, "aws/ec2/i-*\n")
	rs, _ := ignore.LoadFile(path)

	if !rs.Matches("aws", "ec2", "i-aabbcc") {
		t.Error("expected glob pattern to match i-aabbcc")
	}
	if rs.Matches("aws", "ec2", "sg-001") {
		t.Error("expected glob pattern not to match sg-001")
	}
}

func TestDefaultPath(t *testing.T) {
	got := ignore.DefaultPath("/home/user/project")
	want := "/home/user/project/.driftignore"
	if got != want {
		t.Errorf("DefaultPath = %q, want %q", got, want)
	}
}
