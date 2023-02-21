package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/awalterschulze/gographviz"
	"github.com/sinkingpoint/kiora/cmd/kiora/config/nodes"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/kiora/lib/kiora/notify"
)

type ConfigFile struct {
	nodes map[string]nodes.Node
	links map[string][]Link
}

func (c *ConfigFile) GetNotifiersForAlert(a *model.Alert) []notify.Notifier {
	leaves := []notify.Notifier{}

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

		if node, ok := c.nodes[nodeName].(notify.Notifier); node != nil && ok {
			leaves = append(leaves, node)
		}
	}

	return leaves
}

// Validate returns nil if the config is valid, or an error to be displayed to the user if not.
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

func (c *ConfigFile) AmIAuthoritativeFor(a *model.Alert) bool {
	return false
}

// LoadConfigFile reads the given file, and parses it into a config, returning any parsing errors.
func LoadConfigFile(path string) (*ConfigFile, error) {
	conf := &ConfigFile{
		nodes: make(map[string]nodes.Node),
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
		cons := nodes.LookupNode(nodeType)
		if cons == nil {
			return conf, fmt.Errorf("invalid node type: %q", nodeType)
		}

		node, err := cons(rawNode.name, rawNode.attrs)
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
