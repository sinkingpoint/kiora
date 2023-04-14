package config

type Node interface {
	// Type returns a static name of the type of the node for debugging purposes.
	Type() string
}

// NodeConstructor is a function that takes a raw graph node and turns it into a node that can actually process things.
type NodeConstructor = func(name string, bus NodeBus, attrs map[string]string) (Node, error)

var nodeRegistry = map[string]NodeConstructor{}

func RegisterNode(name string, constructor NodeConstructor) {
	nodeRegistry[name] = constructor
}

// LookupNode takes a node type name and returns a constructor that can be used to make nodes of that name.
func LookupNode(name string) (NodeConstructor, bool) {
	cons, ok := nodeRegistry[name]
	return cons, ok
}
