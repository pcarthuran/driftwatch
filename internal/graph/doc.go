// Package graph provides utilities for building and rendering a dependency
// graph of infrastructure resources based on drift detection results.
//
// # Overview
//
// Build accepts a slice of drift.Result values and constructs a Graph whose
// nodes represent individual resources. Resources that share the same
// provider and type are connected by undirected peer edges, making it easy
// to spot clusters of related drift.
//
// # Output format
//
// Write serialises a Graph to the Graphviz DOT language. Drifted nodes are
// coloured salmon; clean nodes are coloured lightblue. The output can be
// piped directly to `dot -Tpng` or any other Graphviz renderer.
//
// Example:
//
//	g := graph.Build(results)
//	graph.Write(g, os.Stdout)
package graph
