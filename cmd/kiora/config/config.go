package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/awalterschulze/gographviz"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type ConfigFile struct {
	nodes map[string]Node
	links map[string][]Link
}

func (c *ConfigFile) GetNotifiersForAlert(a *model.Alert) []kioradb.ModelWriter {
	leaves := []kioradb.ModelWriter{}

	// We expect here that the ConfigFile has been passed through `Validate` already, and thus
	// is assumed to have no cycles.
	stack := []string{"alerts"}
	for len(stack) > 0 {
		nodeName := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, link := range c.links[nodeName] {
			if link.incomingFilter == nil || link.incomingFilter.FilterAlert(a) {
				stack = append(stack, link.to)
			}
		}

		if node, ok := c.nodes[nodeName].(kioradb.ModelWriter); node != nil && ok {
			leaves = append(leaves, node)
		}
	}

	return leaves
}

func (c *ConfigFile) Validate() error {
	// Check if the config file is acyclic.
	for _, tree := range []string{"alerts", "silences"} {
		stack := []string{tree}
		visited := map[string]bool{}
		for len(stack) > 0 {
			nodeName := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if visited[nodeName] {
				return errors.New("config graph cannot contain loops")
			}

			visited[nodeName] = true
			for _, link := range c.links[nodeName] {
				stack = append(stack, link.to)
			}
		}
	}

	return nil
}

// LoadConfigFile reads the given file, and parses it into a config, returning any parsing errors.
func LoadConfigFile(path string) (*ConfigFile, error) {
	conf := &ConfigFile{
		nodes: make(map[string]Node),
		links: make(map[string][]Link),
	}

	body, err := os.ReadFile(path)
	if err != nil {
		return conf, err
	}

	graphAst, err := gographviz.Parse(body)
	if err != nil {
		return conf, err
	}

	configGraph := newConfigGraph()
	if err := gographviz.Analyse(graphAst, &configGraph); err != nil {
		return conf, err
	}

	for _, rawNode := range configGraph.nodes {
		nodeType := rawNode.attrs["type"]
		cons := LookupNode(nodeType)
		if cons == nil {
			return conf, fmt.Errorf("invalid node type: %q", nodeType)
		}

		node, err := cons(rawNode)
		if err != nil {
			return conf, err
		}

		conf.nodes[rawNode.name] = node
	}

	for _, rawLink := range configGraph.edges {
		linkType := rawLink.attrs["type"]
		cons := LookupFilter(linkType)

		var filter Filter
		if cons != nil {
			filter, err = cons(rawLink)
			if err != nil {
				return conf, err
			}
		}

		conf.links[rawLink.from] = append(conf.links[rawLink.from], Link{
			to:             rawLink.to,
			incomingFilter: filter,
		})
	}

	return conf, conf.Validate()
}
