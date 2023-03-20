package config

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// calculateRootsFrom starts at the given node name and walks back up the
// tree, returning the name of all the nodes that have no parents.
func calculateRootsFrom(graph *ConfigFile, nodeName string) HashSet {
	roots := HashSet{}
	visited := HashSet{}

	stack := []string{ACK_LEAF}
	for len(stack) > 0 {
		nodeName := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := visited[nodeName]; ok {
			continue
		}

		visited[nodeName] = struct{}{}

		if len(graph.reverseLinks[nodeName]) == 0 {
			roots[nodeName] = struct{}{}
		} else {
			for _, link := range graph.reverseLinks[nodeName] {
				stack = append(stack, link.to)
			}
		}
	}

	return roots
}

// searchForAckNode starts at the given fromNode and does a depth-first search across the graph,
// checking the filters on each link and trying to find a path to the given destinationNode,
// returning whether or not it was able to find it.
func searchForNode(ctx context.Context, graph *ConfigFile, fromNode string, destinationNode string, ack *model.AlertAcknowledgement) bool {
	if fromNode == destinationNode {
		return true
	}

	for _, link := range graph.links[fromNode] {
		if !link.incomingFilter.Filter(ctx, ack) {
			continue
		}

		if searchForNode(ctx, graph, link.to, destinationNode, ack) {
			return true
		}
	}

	return false
}
