package notify

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Notifier represents something that can send a notification about an alert.
type Notifier interface {
	Notify(ctx context.Context, alerts ...model.Alert) error
}

// NotifierConfig represents a configuration that can return a list of notifiers for a given alert.
type NotifierConfig interface {
	// Returns the notifiers that should be invoked for the given alert. If the response is nil,
	// then the notifier should do nothing, as opposed to an empty array that represents that the alert
	// should be processed as if it should be considered to be properly notified.
	GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []Notifier
}

// ClusterNotifier binds a notifier and a clusterer, returning no notifiers if this node isn't responsible for the alert.
type ClusterNotifier struct {
	clusterer clustering.Clusterer
	conf      NotifierConfig
}

func NewClusterNotifier(clusterer clustering.Clusterer, notifier NotifierConfig) *ClusterNotifier {
	return &ClusterNotifier{
		clusterer: clusterer,
		conf:      notifier,
	}
}

func (c *ClusterNotifier) GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []Notifier {
	if !c.clusterer.IsAuthoritativeFor(ctx, alert) {
		return nil
	}

	return c.conf.GetNotifiersForAlert(ctx, alert)
}
