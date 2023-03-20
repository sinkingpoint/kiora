package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/awalterschulze/gographviz"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/services/notify/notify_config"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

const ALERT_ROOT = "alerts"
const SILENCES_ROOT = "silences"
const ACK_LEAF = "acks"

var _ = notify_config.Config(&ConfigFile{})

// Link represents a connection between nodes, that may or may not have an attached filter.
type Link struct {
	incomingFilter config.Filter
	to             string
}

type ConfigFile struct {
	nodes        map[string]config.Node
	links        map[string][]Link
	reverseLinks map[string][]Link
}

func (c *ConfigFile) GetNotifiersForAlert(ctx context.Context, a *model.Alert) []notify_config.Notifier {
	leaves := []notify_config.Notifier{}

	// We expect here that the ConfigFile has been passed through `Validate` already, and thus
	// is assumed to have no cycles.
	stack := []string{ALERT_ROOT}
	for len(stack) > 0 {
		nodeName := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, link := range c.links[nodeName] {
			matchesFilter := true
			if filter, ok := link.incomingFilter.(config.AlertFilter); ok && filter != nil {
				matchesFilter = filter.FilterAlert(a)
			}

			if link.incomingFilter == nil || matchesFilter {
				stack = append(stack, link.to)
			}
		}

		if node, ok := c.nodes[nodeName].(notify_config.Notifier); node != nil && ok {
			leaves = append(leaves, node)
		}
	}

	return leaves
}

// searchForAckNode returns true if we can follow its link, obeying filters, and get the the `acks` node.
func (c *ConfigFile) searchForAckNode(ctx context.Context, nodeName string, alert *model.Alert, ack *model.AlertAcknowledgement) bool {
	if nodeName == ACK_LEAF {
		return true
	}

	for _, link := range c.links[nodeName] {
		if filter, ok := link.incomingFilter.(config.AlertAcknowledgementFilter); ok {
			if !filter.FilterAlertAcknowledgement(alert, ack) {
				continue
			}
		} else {
			log.Warn().Msgf("filter %q at %q is not an AlertAcknowledgementFilter", link.incomingFilter.Type(), nodeName)
		}

		if c.searchForAckNode(ctx, link.to, alert, ack) {
			return true
		}
	}

	return false
}

// AlertAcknowledgementIsValid returns true if we can find a path to the acks node from the roots of the graph.
func (c *ConfigFile) AlertAcknowledgementIsValid(ctx context.Context, alert *model.Alert, ack *model.AlertAcknowledgement) bool {
	roots := c.calculateAckRoots() // TODO(cdouch): memoize this.
	if len(roots) == 0 {
		return true
	}

	for root := range roots {
		node := root
		if c.searchForAckNode(ctx, node, alert, ack) {
			return true
		}
	}

	return false
}

// calculateAckRoots starts at the `acks` node, and works backwards to find the leaves of the acks tree.
func (c *ConfigFile) calculateAckRoots() HashSet {
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

		if len(c.reverseLinks[nodeName]) == 0 {
			roots[nodeName] = struct{}{}
		} else {
			for _, link := range c.reverseLinks[nodeName] {
				stack = append(stack, link.to)
			}
		}
	}

	return roots
}

// LoadConfigFile reads the given file, and parses it into a config, returning any parsing errors.
func LoadConfigFile(path string) (*ConfigFile, error) {
	conf := &ConfigFile{
		nodes:        make(map[string]config.Node),
		links:        make(map[string][]Link),
		reverseLinks: make(map[string][]Link),
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
		cons, ok := config.LookupNode(nodeType)
		if !ok {
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
		cons, ok := config.LookupFilter(linkType)

		if !ok {
			return conf, fmt.Errorf("invalid link type: %q", linkType)
		}

		filter, err := cons(rawLink.attrs)
		if err != nil {
			return conf, err
		}

		if filter == nil {
			panic(fmt.Sprintf("BUG: filter %q produced a nil filter", linkType))
		}

		conf.links[rawLink.from] = append(conf.links[rawLink.from], Link{
			to:             rawLink.to,
			incomingFilter: filter,
		})

		conf.reverseLinks[rawLink.to] = append(conf.reverseLinks[rawLink.to], Link{
			to:             rawLink.from,
			incomingFilter: filter,
		})
	}

	return conf, conf.Validate()
}

// validateConfIsAcyclic starts at the given roots, and validates that there are no cycles in the
// graph when starting at them. This makes sure that we don't get infinite notification loops.
func (c *ConfigFile) validateConfIsAcyclic(roots map[string]struct{}) error {
	for tree := range roots {
		// Construct a stack, and a set of visited nodes. We'll use the stack to do a DFS, and
		// the set to make sure we don't visit the same node twice.
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

// validateConfHasValidRoots makes sure that the config doesn't have any links going into the
// nodes that are expected to be roots. Roots are expected to be entrypoints into the config,
// like the POST alerts API, so things going into them doesn't make sense.
func (c *ConfigFile) validateConfHasValidRoots(roots HashSet) error {
	for _, link := range c.links {
		for _, l := range link {
			if _, ok := roots[l.to]; ok {
				// We have a link into a root node, which is not allowed.
				return fmt.Errorf("invalid link going into root node: %q", l.to)
			}
		}
	}

	return nil
}

// validateConfHasValidLeaves makes sure that the config doesn't have any links going out of the
// nodes that are expected to be leaves. Leaves are expected to be the end of the config, so
// things going out of them doesn't make sense.
func (c *ConfigFile) validateConfHasValidLeaves(leaves HashSet) error {
	for leaf := range leaves {
		if linksFrom := c.links[leaf]; len(linksFrom) > 0 {
			// We have a link from a leaf node, which is not allowed.
			return fmt.Errorf("invalid link going from leaf node: %q", leaf)
		}
	}
	return nil
}

// Validate returns nil if the config is valid, or an error to be displayed to the user if not.
func (c *ConfigFile) Validate() error {
	roots := toHashSet([]string{ALERT_ROOT, SILENCES_ROOT})
	leaves := toHashSet([]string{ACK_LEAF})

	if err := c.validateConfIsAcyclic(roots); err != nil {
		return err
	}

	if err := c.validateConfHasValidRoots(roots); err != nil {
		return err
	}

	if err := c.validateConfHasValidLeaves(leaves); err != nil {
		return err
	}

	return nil
}

// HashSet is a helper type that manages a set of strings.
type HashSet map[string]struct{}

func toHashSet(s []string) HashSet {
	set := HashSet{}
	for _, v := range s {
		set[v] = struct{}{}
	}

	return set
}
