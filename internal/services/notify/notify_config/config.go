package notify_config

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// ClusterNotifier binds a notifier and a clusterer, returning no notifiers if this node isn't responsible for the alert.
type ClusterNotifier struct {
	clusterer clustering.Clusterer
	conf      config.Config
}

func NewClusterNotifier(clusterer clustering.Clusterer, notifier config.Config) *ClusterNotifier {
	return &ClusterNotifier{
		clusterer: clusterer,
		conf:      notifier,
	}
}

func (c *ClusterNotifier) GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []config.Notifier {
	if !c.clusterer.IsAuthoritativeFor(ctx, alert) {
		return nil
	}

	return c.conf.GetNotifiersForAlert(ctx, alert)
}
