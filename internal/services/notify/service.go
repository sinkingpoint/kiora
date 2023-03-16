package notify

import (
	"context"
	"time"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/server/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ = services.Service(&NotifyService{})

// NotifyService is a background service that scans the db for alerts to send notifications for.
type NotifyService struct {
	config      NotifierConfig
	db          kioradb.DB
	broadcaster clustering.Broadcaster
}

func NewNotifyService(config NotifierConfig, db kioradb.DB, broadcaster clustering.Broadcaster) *NotifyService {
	return &NotifyService{
		config:      config,
		db:          db,
		broadcaster: broadcaster,
	}
}

func (n *NotifyService) Name() string {
	return "notify"
}

func (n *NotifyService) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
outer:
	for {
		select {
		case <-ticker.C:
			n.notify(ctx)
		case <-ctx.Done():
			break outer
		}
	}
	return nil
}

func (n *NotifyService) notify(ctx context.Context) {
	q := &query.StatusQuery{
		Status: model.AlertStatusProcessing,
	}

	for _, a := range n.db.QueryAlerts(ctx, q) {
		n.notifyAlert(ctx, a)
	}
}

// notifyAlert sends a notification for the given alert.
// TODO(cdouch): Handle errors here.
func (n *NotifyService) notifyAlert(ctx context.Context, a model.Alert) {
	notifiers := n.config.GetNotifiersForAlert(ctx, &a)
	if notifiers == nil {
		return
	}

	a.Status = model.AlertStatusFiring

	for _, n := range notifiers {
		n.Notify(ctx, a) //nolint
	}

	n.broadcaster.BroadcastAlerts(ctx, a) //nolint
}
