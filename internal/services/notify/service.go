package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/server/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ = services.Service(&NotifyService{})

const DEFAULT_RENOTIFY_INTERVAL = 3 * time.Hour

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
			n.notifyFiring(ctx)
			n.notifyResolved(ctx)
		case <-ctx.Done():
			break outer
		}
	}
	return nil
}

func (n *NotifyService) notifyFiring(ctx context.Context) {
	q := query.All(query.Status(model.AlertStatusFiring), query.LastNotifyTimeMax(time.Now().Add(-DEFAULT_RENOTIFY_INTERVAL)))

	for _, a := range n.db.QueryAlerts(ctx, q) {
		n.notifyAlert(ctx, a)
	}
}

func (n *NotifyService) notifyResolved(ctx context.Context) {
	q := query.Status(model.AlertStatusResolved)
	for _, alert := range n.db.QueryAlerts(ctx, q) {
		if alert.LastNotifyTime.Before(alert.EndTime) {
			n.notifyAlert(ctx, alert)
		}
	}
}

// notifyAlert sends a notification for the given alert.
// TODO(cdouch): Handle errors here.
func (n *NotifyService) notifyAlert(ctx context.Context, a model.Alert) {
	ctx, span := otel.Tracer("").Start(ctx, "NotifyService.notifyAlert")
	defer span.End()

	span.SetAttributes(attribute.String("alert", fmt.Sprintf("%+v", a)))

	notifiers := n.config.GetNotifiersForAlert(ctx, &a)
	if notifiers == nil {
		span.AddEvent("Not responsible for this alert")
		return
	}

	a.LastNotifyTime = time.Now()

	for _, n := range notifiers {
		n.Notify(ctx, a) //nolint
	}

	n.broadcaster.BroadcastAlerts(ctx, a) //nolint
}
