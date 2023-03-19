package clustering

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Broadcaster defines something that can tell other things about data.
type Broadcaster interface {
	// BroadcastAlerts broadcasts a group of alerts to a cluster.
	BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error

	// BroadcastAlertAcknowledgement broadcasts an AlertAcknowledgement of the given alert.
	BroadcastAlertAcknowledgement(ctx context.Context, alertID string, ack model.AlertAcknowledgement) error
}
