package nodes

// AnchorNode is the default node type, if nothing else is specified. They do nothing except
// act as anchor points for Links to allow splitting one or more incoming links into one or more outgoing ones.
type AnchorNode struct {
}

func (a *AnchorNode) Type() string {
	return "anchor"
}
