package clustering

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Clusterer is used to determine if this node is authoritative (and thus should send a notification for) a given alert.
type Clusterer interface {
	// IsAuthoritativeFor returns true if this is the node that should send notifications for the given alert.
	IsAuthoritativeFor(ctx context.Context, a *model.Alert) bool

	// Nodes returns a list of the nodes in the cluster.
	Nodes() []any
}

// ClustererDelegates receive cluster updates (node additions and removals).
type ClustererDelegate interface {
	// AddNode is called when a node is added to the cluster.
	AddNode(name, address string)

	// RemoveNode is called when a node is removed, or fails.
	RemoveNode(name string)
}

// EventDelegate provides a delegate that can handle events as they come in from a cluster channel.
type EventDelegate interface {
	// ProcessAlert is called when a new alert comes in. There are no guarantees that this alert isn't one
	// we haven't seen before - it might be an update on status etc.
	ProcessAlert(ctx context.Context, alert model.Alert)

	// ProcessAlertAcknowledgement is called when a new alert acknowledgement comes in.
	ProcessAlertAcknowledgement(ctx context.Context, alertID string, ack model.AlertAcknowledgement)

	// ProcessSilence is called when a new silence comes in. There are no guarantees that this silence isn't one
	// we haven't seen before - it might be an update on status etc.
	ProcessSilence(ctx context.Context, silence model.Silence)
}
