package nodes

type Node interface {
	// Type returns a static name of the type of the node for debugging purposes.
	Type() string
}

// NodeConstructor is a function that takes a raw graph node and turns it into a node that can actually process things.
type NodeConstructor = func(name string, attrs map[string]string) (Node, error)

var nodeRegistry = map[string]NodeConstructor{
	"":       func(name string, attrs map[string]string) (Node, error) { return &AnchorNode{}, nil },
	"stdout": NewFileNotifierNode,
	"stderr": NewFileNotifierNode,
	"file":   NewFileNotifierNode,
}

// LookupNode takes a node type name and returns a constructor that can be used to make nodes of that name.
func LookupNode(name string) NodeConstructor {
	return nodeRegistry[name]
}
