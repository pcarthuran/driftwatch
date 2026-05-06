package graph_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/graph"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "res-1", Provider: "aws", ResourceType: "s3", Status: drift.StatusClean},
		{ResourceID: "res-2", Provider: "aws", ResourceType: "s3", Status: drift.StatusModified},
		{ResourceID: "res-3", Provider: "gcp", ResourceType: "gcs", Status: drift.StatusMissing},
	}
}

func TestBuild_NodeCount(t *testing.T) {
	g := graph.Build(sampleResults())
	if len(g.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.Nodes))
	}
}

func TestBuild_DriftedFlag(t *testing.T) {
	g := graph.Build(sampleResults())
	if g.Nodes["res-1"].Drifted {
		t.Error("res-1 should not be drifted")
	}
	if !g.Nodes["res-2"].Drifted {
		t.Error("res-2 should be drifted")
	}
	if !g.Nodes["res-3"].Drifted {
		t.Error("res-3 should be drifted")
	}
}

func TestBuild_PeerEdges(t *testing.T) {
	g := graph.Build(sampleResults())
	// res-1 and res-2 are both aws/s3, so they should be linked.
	edges := g.Nodes["res-1"].Edges
	if len(edges) != 1 || edges[0] != "res-2" {
		t.Errorf("expected res-1 to have edge to res-2, got %v", edges)
	}
}

func TestBuild_NoEdgesAcrossTypes(t *testing.T) {
	g := graph.Build(sampleResults())
	// res-3 is gcp/gcs — no peers.
	if len(g.Nodes["res-3"].Edges) != 0 {
		t.Errorf("res-3 should have no edges, got %v", g.Nodes["res-3"].Edges)
	}
}

func TestWrite_DOTContainsNodes(t *testing.T) {
	g := graph.Build(sampleResults())
	var buf bytes.Buffer
	if err := graph.Write(g, &buf); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, id := range []string{"res-1", "res-2", "res-3"} {
		if !strings.Contains(out, id) {
			t.Errorf("output missing node %q", id)
		}
	}
}

func TestWrite_DOTContainsDigraph(t *testing.T) {
	g := graph.Build(sampleResults())
	var buf bytes.Buffer
	_ = graph.Write(g, &buf)
	if !strings.HasPrefix(buf.String(), "digraph drift {") {
		t.Error("output should start with 'digraph drift {'")
	}
}

func TestWrite_EmptyGraph(t *testing.T) {
	g := graph.Build(nil)
	var buf bytes.Buffer
	if err := graph.Write(g, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "digraph drift") {
		t.Error("empty graph should still produce valid DOT header")
	}
}
