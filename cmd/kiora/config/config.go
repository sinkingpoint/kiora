package config

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/awalterschulze/gographviz"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel"
)

const (
	ALERT_ROOT    = "alerts"
	SILENCES_LEAF = "silences"
	ACK_LEAF      = "acks"
)

var _ = config.Config(&ConfigFile{})

type globalOptions struct {
	TenantKey *template.Template `config:"tenant_key"`
}

// Link represents a connection between nodes, that may or may not have an attached filter.
type Link struct {
	incomingFilter config.Filter
	to             string
}

type ConfigFile struct {
	// nodes is a map of node name to node. This is used to look up nodes by name.
	nodes map[string]config.Node

	// links is a map of node name to a list of links that go out of that node.
	links map[string][]Link

	// reverseLinks is a map of node name to a list of links that go into that node, for backwards traversal.
	reverseLinks map[string][]Link
}

// GetNotifiersForAlert walks the config graph, building up notification settings as we go before returning a list
// of notifiers we hit along the way. We expect here that the ConfigFile has been passed through `Validate` already, and thus
// is assumed to have no cycles.
func (c *ConfigFile) GetNotifiersForAlert(ctx context.Context, a *model.Alert) []config.NotifierSettings {
	ctx, span := otel.Tracer("").Start(ctx, "ConfigFile.GetNotifiersForAlert")
	defer span.End()

	leaves := []config.NotifierSettings{}

	// nodeMeta is a node that we've traversed to, and the partial configuration that we've built up along the path there.
	// TODO(cdouch): I'm not _entirely_ sure what happens when we get a two paths to the same node, but with different
	// configurations. I think we'll end up with two notifiers, but I'm not sure. Need to think about this more.
	type nodeMeta struct {
		name        string
		partialConf config.NotifierSettings
	}

	// We use a stack here to do a depth-first search of the graph, starting at the `alerts` node.
	stack := []nodeMeta{{
		name:        ALERT_ROOT,
		partialConf: config.DefaultNotifierSettings(),
	}}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if confNode, ok := c.nodes[node.name].(config.NotifierSettingsNode); confNode != nil && ok {
			if err := confNode.Apply(&node.partialConf); err != nil {
				log.Warn().Err(err).Msg("failed to apply notifier settings node")
			}
		}

		for _, link := range c.links[node.name] {
			matchesFilter := link.incomingFilter == nil || link.incomingFilter.Filter(ctx, a) != nil
			if matchesFilter {
				stack = append(stack, nodeMeta{
					name:        link.to,
					partialConf: node.partialConf,
				})
			}
		}

		if notifier, ok := c.nodes[node.name].(config.Notifier); notifier != nil && ok {
			leaves = append(leaves, node.partialConf.WithNotifier(notifier))
		}
	}

	return leaves
}

// validateData walks the config graph, along every path into the given leaf. We check every path against the given Fielder,
// and return an error if we can't find a path into the leaf that matches the data.
func (c *ConfigFile) validateData(ctx context.Context, leaf string, data config.Fielder) error {
	roots := calculateRootsFrom(c, leaf) // TODO(cdouch): memoize this.
	if len(roots) == 0 {
		return nil
	}

	var allErrs error
	for root := range roots {
		if err := searchForNode(ctx, c, root, leaf, data); err == nil {
			return nil
		} else {
			allErrs = multierror.Append(allErrs, err)
		}
	}

	return allErrs
}

func (c *ConfigFile) ValidateData(ctx context.Context, data config.Fielder) error {
	switch data.(type) {
	case *model.AlertAcknowledgement:
		return c.validateData(ctx, ACK_LEAF, data)
	case *model.Silence:
		return c.validateData(ctx, SILENCES_LEAF, data)
	default:
		panic("BUG: unhandled data validation")
	}
}

// LoadConfigFile reads the given file, and parses it into a config, returning any parsing errors.
func LoadConfigFile(path string, logger zerolog.Logger) (*ConfigFile, error) {
	conf := &ConfigFile{
		nodes:        make(map[string]config.Node),
		links:        make(map[string][]Link),
		reverseLinks: make(map[string][]Link),
	}

	body, err := os.ReadFile(path)
	if err != nil {
		return conf, errors.Wrap(err, "failed to read config file")
	}

	graphAst, err := gographviz.Parse(body)
	if err != nil {
		return conf, errors.Wrap(err, "failed to parse config file as dot")
	}

	configGraph := newConfigGraph()
	if err := gographviz.Analyse(graphAst, &configGraph); err != nil {
		return conf, errors.Wrap(err, "failed to load config file")
	}

	options := globalOptions{}

	if err := unmarshal.UnmarshalConfig(configGraph.attrs, &options, unmarshal.UnmarshalOpts{DisallowUnknownFields: true}); err != nil {
		return conf, errors.Wrap(err, "failed to parse config file")
	}

	var tenanter config.Tenanter
	if options.TenantKey != nil {
		tenanter = config.NewTemplateTenanter(options.TenantKey)
	}

	globals := config.NewGlobals(config.WithLogger(logger), config.WithTenanter(tenanter))

	for _, rawNode := range configGraph.nodes {
		nodeType := rawNode.attrs["type"]
		cons, ok := config.LookupNode(nodeType)
		if !ok {
			return conf, fmt.Errorf("invalid node type: %q", nodeType)
		}

		node, err := cons(rawNode.name, globals, rawNode.attrs)
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

		filter, err := cons(globals, rawLink.attrs)
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
	roots := toHashSet([]string{ALERT_ROOT})
	leaves := toHashSet([]string{ACK_LEAF, SILENCES_LEAF})

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
