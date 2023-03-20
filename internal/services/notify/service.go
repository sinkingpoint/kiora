package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/internal/services/notify/notify_config"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ = services.Service(&NotifyService{})

const DEFAULT_RENOTIFY_INTERVAL = 3 * time.Hour

// NotifyService is a background service that scans the db for alerts to send notifications for.
type NotifyService struct {
	config notify_config.Config
	bus    services.Bus
}

func NewNotifyService(config notify_config.Config, bus services.Bus) *NotifyService {
	return &NotifyService{
		config: config,
		bus:    bus,
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
	q := query.AllAlerts(query.Status(model.AlertStatusFiring), query.LastNotifyTimeMax(stubs.Time.Now().Add(-DEFAULT_RENOTIFY_INTERVAL)))

	for _, a := range n.bus.DB().QueryAlerts(ctx, q) {
		n.notifyAlert(ctx, a)
	}
}

func (n *NotifyService) notifyResolved(ctx context.Context) {
	q := query.Status(model.AlertStatusResolved)
	for _, alert := range n.bus.DB().QueryAlerts(ctx, q) {
		if alert.LastNotifyTime.Before(alert.EndTime) {
			n.notifyAlert(ctx, alert)
		}
	}
}

// notifyAlert sends a notification for the given alert.
func (n *NotifyService) notifyAlert(ctx context.Context, a model.Alert) {
	ctx, span := otel.Tracer("").Start(ctx, "NotifyService.notifyAlert")
	defer span.End()

	span.SetAttributes(attribute.String("alert", fmt.Sprintf("%+v", a)))

	notifiers := n.config.GetNotifiersForAlert(ctx, &a)
	if notifiers == nil {
		span.AddEvent("Not responsible for this alert")
		return
	}

	a.LastNotifyTime = stubs.Time.Now()

	for _, notifier := range notifiers {
		if err := notifier.Notify(ctx, a); err != nil {
			n.bus.Logger("notify").Err(err).Msg("failed to notify for alert")
		}
	}

	if err := n.bus.Broadcaster().BroadcastAlerts(ctx, a); err != nil {
		n.bus.Logger("notify").Err(err).Msg("failed to broadcast the sucessful notify")
	}
}
