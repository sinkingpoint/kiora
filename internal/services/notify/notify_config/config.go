package notify_config

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// ClusterNotifier binds a notifier and a clusterer, returning no notifiers if this node isn't responsible for the alert.
type ClusterNotifier struct {
	config.Config
	clusterer clustering.Clusterer
}

func NewClusterNotifier(clusterer clustering.Clusterer, config config.Config) *ClusterNotifier {
	return &ClusterNotifier{
		Config:    config,
		clusterer: clusterer,
	}
}

func (c *ClusterNotifier) GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []config.Notifier {
	if !c.clusterer.IsAuthoritativeFor(ctx, alert) {
		return nil
	}

	return c.Config.GetNotifiersForAlert(ctx, alert)
}
