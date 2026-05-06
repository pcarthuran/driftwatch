// Package graph builds a dependency graph from drift detection results,
// allowing callers to visualise relationships between drifted resources.
package graph

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Node represents a single resource vertex in the graph.
type Node struct {
	ID       string
	Provider string
	Type     string
	Drifted  bool
	Edges    []string // IDs of nodes this node depends on
}

// Graph holds all nodes keyed by resource ID.
type Graph struct {
	Nodes map[string]*Node
}

// Build constructs a Graph from drift detection results.
// Resources that share the same provider+type are connected by implicit edges.
func Build(results []drift.Result) *Graph {
	g := &Graph{Nodes: make(map[string]*Node)}

	for _, r := range results {
		n := &Node{
			ID:       r.ResourceID,
			Provider: r.Provider,
			Type:     r.ResourceType,
			Drifted:  r.Status != drift.StatusClean,
		}
		g.Nodes[r.ResourceID] = n
	}

	// Link nodes of the same provider+type as peers.
	groups := make(map[string][]string)
	for id, n := range g.Nodes {
		key := n.Provider + "/" + n.Type
		groups[key] = append(groups[key], id)
	}
	for _, ids := range groups {
		sort.Strings(ids)
		for i, id := range ids {
			for j, peer := range ids {
				if i != j {
					g.Nodes[id].Edges = append(g.Nodes[id].Edges, peer)
				}
			}
		}
	}

	return g
}

// Write renders the graph as a simple DOT-format string to w.
func Write(g *Graph, w io.Writer) error {
	fmt.Fprintln(w, "digraph drift {")
	fmt.Fprintln(w, "  rankdir=LR;")

	keys := make([]string, 0, len(g.Nodes))
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, id := range keys {
		n := g.Nodes[id]
		color := "lightblue"
		if n.Drifted {
			color = "salmon"
		}
		label := fmt.Sprintf("%s\\n%s/%s", id, n.Provider, n.Type)
		fmt.Fprintf(w, "  %q [label=%q style=filled fillcolor=%s];\n",
			id, label, color)
	}

	seen := make(map[string]bool)
	for _, id := range keys {
		n := g.Nodes[id]
		for _, peer := range n.Edges {
			edge := id + "->" + peer
			reverse := peer + "->" + id
			if seen[edge] || seen[reverse] {
				continue
			}
			seen[edge] = true
			fmt.Fprintf(w, "  %q -> %q [dir=none];\n", id, peer)
		}
	}

	fmt.Fprintln(w, "}")
	_ = strings.TrimSpace // keep import
	return nil
}
