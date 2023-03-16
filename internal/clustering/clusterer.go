package clustering

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Clusterer is used to determine if this node is authoritative (and thus should send a notification for) a given alert.
type Clusterer interface {
	IsAuthoritativeFor(ctx context.Context, a *model.Alert) bool
}

// ClustererDelegates receive cluster updates (node additions and removals)
type ClustererDelegate interface {
	// AddNode is called when a node is added to the cluster.
	AddNode(name string, address string)

	// RemoveNode is called when a node is removed, or fails.
	RemoveNode(name string)
}
