package graph_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/graph"
)

func TestIntegration_AllDrifted_AllSalmon(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "a", Provider: "aws", ResourceType: "ec2", Status: drift.StatusModified},
		{ResourceID: "b", Provider: "aws", ResourceType: "ec2", Status: drift.StatusExtra},
	}
	g := graph.Build(results)
	var buf bytes.Buffer
	_ = graph.Write(g, &buf)
	out := buf.String()
	count := strings.Count(out, "salmon")
	if count != 2 {
		t.Errorf("expected 2 salmon nodes, got %d", count)
	}
}

func TestIntegration_MixedProviders_EdgeIsolation(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "x1", Provider: "aws", ResourceType: "rds", Status: drift.StatusClean},
		{ResourceID: "x2", Provider: "azure", ResourceType: "rds", Status: drift.StatusClean},
	}
	// Same type but different providers — should NOT share edges.
	g := graph.Build(results)
	if len(g.Nodes["x1"].Edges) != 0 {
		t.Errorf("x1 should have no edges across providers, got %v", g.Nodes["x1"].Edges)
	}
}

func TestIntegration_SingleResource_NoEdges(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "solo", Provider: "gcp", ResourceType: "bigquery", Status: drift.StatusClean},
	}
	g := graph.Build(results)
	if len(g.Nodes["solo"].Edges) != 0 {
		t.Errorf("solo node should have no edges, got %v", g.Nodes["solo"].Edges)
	}
}
