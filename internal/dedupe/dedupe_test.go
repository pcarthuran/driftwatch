package dedupe_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/dedupe"
	"github.com/driftwatch/internal/drift"
)

func sampleResults(drifted bool) []drift.Result {
	return []drift.Result{
		{Provider: "aws", Type: "ec2", ID: "i-001", Drifted: drifted},
		{Provider: "gcp", Type: "vm", ID: "vm-a", Drifted: false},
	}
}

func TestDeduplicate_NoDuplicates(t *testing.T) {
	input := sampleResults(false)
	out := dedupe.Deduplicate(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestDeduplicate_RemovesDuplicates(t *testing.T) {
	a := []drift.Result{{Provider: "aws", Type: "ec2", ID: "i-001", Drifted: false}}
	b := []drift.Result{{Provider: "aws", Type: "ec2", ID: "i-001", Drifted: true}}

	out := dedupe.Deduplicate(a, b)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	// last-seen wins
	if !out[0].Drifted {
		t.Error("expected last-seen (drifted=true) to win")
	}
}

func TestDeduplicate_MergesMultipleSets(t *testing.T) {
	a := sampleResults(false)
	b := []drift.Result{
		{Provider: "azure", Type: "vm", ID: "vm-1", Drifted: true},
	}
	out := dedupe.Deduplicate(a, b)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestDeduplicate_SortedOutput(t *testing.T) {
	input := []drift.Result{
		{Provider: "gcp", Type: "vm", ID: "z"},
		{Provider: "aws", Type: "ec2", ID: "a"},
	}
	out := dedupe.Deduplicate(input)
	if out[0].Provider != "aws" {
		t.Errorf("expected aws first, got %s", out[0].Provider)
	}
}

func TestDeduplicate_Empty(t *testing.T) {
	out := dedupe.Deduplicate()
	if len(out) != 0 {
		t.Errorf("expected empty, got %d", len(out))
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	dedupe.Write(&buf, sampleResults(true))
	out := buf.String()
	if !strings.Contains(out, "PROVIDER") {
		t.Error("expected PROVIDER header")
	}
	if !strings.Contains(out, "drifted") {
		t.Error("expected drifted status")
	}
}

func TestWrite_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	dedupe.Write(&buf, []drift.Result{})
	out := buf.String()
	if !strings.Contains(out, "PROVIDER") {
		t.Error("expected header even with no results")
	}
}
