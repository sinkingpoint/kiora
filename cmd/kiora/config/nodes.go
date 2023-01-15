package config

type Node interface {
	Type() string
}

type NodeConstructor = func(n node) (Node, error)

var nodeRegistry = map[string]NodeConstructor{
	"": func(n node) (Node, error) { return &AnchorNode{}, nil },
}

func LookupNode(name string) NodeConstructor {
	return nodeRegistry[name]
}

// AnchorNode is the default node type, if nothing else is specified. They do nothing except
// act as anchor points for Links to allow splitting one or more incoming links into one or more outgoing ones.
type AnchorNode struct{}

func (a *AnchorNode) Type() string {
	return "anchor"
}
