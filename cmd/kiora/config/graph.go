package config

// Defines a custom Graph for gographviz that accepts all extra attributes.

import (
	"errors"
	"fmt"
	"strings"
)

// node is a node in the config graph that defines a filter, or a receiver.
type node struct {
	name  string
	attrs map[string]string
}

// edge defines an edge between two nodes in the graph.
type edge struct {
	from  string
	to    string
	attrs map[string]string
}

// configGraph is the raw graphviz graph, as loaded from the config file.
type configGraph struct {
	name      string
	attrs     map[string]string
	subGraphs map[string]configGraph
	nodes     map[string]node
	edges     []edge
}

// newConfigGraph constructs a new configGraph, initializing all the maps.
func newConfigGraph() configGraph {
	return configGraph{
		name:      "",
		attrs:     make(map[string]string),
		subGraphs: make(map[string]configGraph),
		nodes:     make(map[string]node),
		edges:     []edge{},
	}
}

func (c *configGraph) SetStrict(strict bool) error {
	return nil
}

func (c *configGraph) SetDir(directed bool) error {
	return nil
}

func (c *configGraph) SetName(name string) error {
	c.name = name
	return nil
}

func (c *configGraph) AddPortEdge(src, srcPort, dst, dstPort string, directed bool, attrs map[string]string) error {
	return c.AddEdge(src, dst, directed, attrs)
}

func (c *configGraph) AddEdge(src, dst string, directed bool, attrs map[string]string) error {
	if !directed {
		return errors.New("edges in the Config Graph must be directed")
	}

	c.edges = append(c.edges, edge{
		from:  src,
		to:    dst,
		attrs: attrs,
	})

	return nil
}

func (c *configGraph) AddNode(parentGraph string, name string, attrs map[string]string) error {
	if parentGraph == c.name {
		if _, ok := c.nodes[name]; ok {
			return fmt.Errorf("config graph already contains a node called %q", name)
		}

		for i := range attrs {
			attrs[i] = strings.Trim(attrs[i], "\"")
		}

		c.nodes[name] = node{
			name,
			attrs,
		}

		return nil
	} else {
		if sub, ok := c.subGraphs[parentGraph]; ok {
			return sub.AddNode(parentGraph, name, attrs)
		} else {
			return fmt.Errorf("failed to find subgraph %q to add node", parentGraph)
		}
	}
}

func (c *configGraph) AddAttr(parentGraph string, field, value string) error {
	if parentGraph == c.name {
		if _, ok := c.attrs[field]; ok {
			return fmt.Errorf("graph already has an attribute %q", field)
		}

		value := strings.Trim(value, "\"")

		c.attrs[field] = value
		return nil
	} else {
		if sub, ok := c.subGraphs[parentGraph]; ok {
			return sub.AddAttr(parentGraph, field, value)
		} else {
			return fmt.Errorf("failed to find subgraph %q to add node", parentGraph)
		}
	}
}

func (c *configGraph) AddSubGraph(parentGraph string, name string, attrs map[string]string) error {
	if parentGraph == c.name {
		if _, ok := c.attrs[name]; ok {
			return fmt.Errorf("graph already has an subgraph %q", name)
		}

		graph := newConfigGraph()
		graph.name = name
		graph.attrs = attrs

		c.subGraphs[name] = graph
		return nil
	} else {
		// We only support one level of nesting for now, error if we're trying to add a subgraph to a subgraph.
		return errors.New("config only supports one layer of nesting")
	}
}

func (c *configGraph) String() string {
	return ""
}
