package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
)

func baseSnapshot(resources []state.Resource) *state.Snapshot {
	return &state.Snapshot{Resources: resources}
}

func TestDetect_NoDrift(t *testing.T) {
	res := []state.Resource{
		{ID: "vm-1", Type: "vm", Fields: map[string]interface{}{"cpu": 2, "mem": "4Gi"}},
	}
	report, err := drift.Detect(baseSnapshot(res), baseSnapshot(res))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Drifted {
		t.Errorf("expected no drift, got %+v", report.Results)
	}
}

func TestDetect_MissingResource(t *testing.T) {
	declared := baseSnapshot([]state.Resource{
		{ID: "vm-1", Type: "vm", Fields: map[string]interface{}{"cpu": 2}},
	})
	live := baseSnapshot([]state.Resource{})

	report, err := drift.Detect(declared, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Drifted {
		t.Fatal("expected drift to be detected")
	}
	if len(report.Results) != 1 || report.Results[0].Status != drift.StatusMissing {
		t.Errorf("expected one missing result, got %+v", report.Results)
	}
}

func TestDetect_ExtraResource(t *testing.T) {
	declared := baseSnapshot([]state.Resource{})
	live := baseSnapshot([]state.Resource{
		{ID: "vm-99", Type: "vm", Fields: map[string]interface{}{"cpu": 4}},
	})

	report, err := drift.Detect(declared, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Drifted {
		t.Fatal("expected drift to be detected")
	}
	if report.Results[0].Status != drift.StatusExtra {
		t.Errorf("expected extra status, got %s", report.Results[0].Status)
	}
}

func TestDetect_ModifiedField(t *testing.T) {
	declared := baseSnapshot([]state.Resource{
		{ID: "db-1", Type: "database", Fields: map[string]interface{}{"size": "10Gi"}},
	})
	live := baseSnapshot([]state.Resource{
		{ID: "db-1", Type: "database", Fields: map[string]interface{}{"size": "20Gi"}},
	})

	report, err := drift.Detect(declared, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Drifted {
		t.Fatal("expected drift to be detected")
	}
	if report.Results[0].Status != drift.StatusModified {
		t.Errorf("expected modified status, got %s", report.Results[0].Status)
	}
	if report.Results[0].Field != "size" {
		t.Errorf("expected field 'size', got %s", report.Results[0].Field)
	}
}

func TestDetect_NilDeclared(t *testing.T) {
	_, err := drift.Detect(nil, baseSnapshot([]state.Resource{}))
	if err == nil {
		t.Error("expected error for nil declared snapshot")
	}
}

func TestDetect_NilLive(t *testing.T) {
	_, err := drift.Detect(baseSnapshot([]state.Resource{}), nil)
	if err == nil {
		t.Error("expected error for nil live snapshot")
	}
}
