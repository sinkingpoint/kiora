package config

import (
	"fmt"
	"os"

	"github.com/awalterschulze/gographviz"
)

type ConfigFile struct {
	nodes map[string]Node
	links map[string][]Link
}

// LoadConfigFile reads the given file, and parses it into a config, returning any parsing errors.
func LoadConfigFile(path string) (ConfigFile, error) {
	conf := ConfigFile{
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

	return conf, nil
}
