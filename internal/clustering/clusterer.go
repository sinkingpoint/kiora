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

// EventDelegate provides a delegate that can handle events as they come in from a cluster channel.
type EventDelegate interface {
	// ProcessAlert is called when a new alert comes in. There are no guarantees that this alert isn't one
	// we haven't seen before - it might be an update on status etc.
	ProcessAlert(ctx context.Context, alert model.Alert)
}
