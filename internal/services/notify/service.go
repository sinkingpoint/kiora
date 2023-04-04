package notify

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ = services.Service(&NotifyService{})

const DEFAULT_RENOTIFY_INTERVAL = 3 * time.Hour

// groupMeta is a helper struct to track groups of alerts that should be notified together.
type groupMeta struct {
	// The group key is a unique identifier for the group of alerts.
	GroupLabels model.Labels

	// The group timeout is the time at which we should send a notification for this group.
	Timeout time.Time

	// The group notifier is the notifier that we should use to send this notification.
	Notifier config.NotifierSettings

	// The group alerts is the list of alerts that are in this group that will be notifier when the timeout is reached.
	Alerts []model.Alert
}

// NotifyService is a background service that scans the db for alerts to send notifications for.
type NotifyService struct {
	config config.Config
	bus    services.Bus

	groupMutex sync.Mutex

	// pendingGroups is a map of notifier names to a list of groups that are pending notification for that notifier.
	pendingGroups map[config.NotifierName][]groupMeta
}

func NewNotifyService(conf config.Config, bus services.Bus) *NotifyService {
	return &NotifyService{
		config:        conf,
		bus:           bus,
		pendingGroups: make(map[config.NotifierName][]groupMeta),
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
			n.notifyGroup(ctx)
		case <-ctx.Done():
			break outer
		}
	}
	return nil
}

func (n *NotifyService) notifyFiring(ctx context.Context) {
	q := query.NewAlertQuery(query.AllAlerts(query.Status(model.AlertStatusFiring), query.LastNotifyTimeMax(stubs.Time.Now().Add(-DEFAULT_RENOTIFY_INTERVAL))))

	for _, a := range n.bus.DB().QueryAlerts(ctx, q) {
		n.notifyAlert(ctx, a)
	}
}

func (n *NotifyService) notifyResolved(ctx context.Context) {
	q := query.Status(model.AlertStatusResolved)
	for _, alert := range n.bus.DB().QueryAlerts(ctx, query.NewAlertQuery(q)) {
		if alert.LastNotifyTime.Before(alert.EndTime) {
			n.notifyAlert(ctx, alert)
		}
	}
}

// notifyGroup will notify all groups that have timed out. This locks the groupMutex for the duration of the function
// which is super expensive and will block all other notifications. This is fine for now, but we should probably
// find a better way to do this.
func (n *NotifyService) notifyGroup(ctx context.Context) {
	n.groupMutex.Lock()
	defer n.groupMutex.Unlock()

	for key, groups := range n.pendingGroups {
		stillWaitingGroups := []groupMeta{}
		for _, g := range groups {
			if g.Timeout.Before(stubs.Time.Now()) {
				if err := g.Notifier.Notify(ctx, g.Alerts...); err != nil {
					n.bus.Logger("notify").Err(err).Msg("failed to notify for alert")
				}
			} else {
				stillWaitingGroups = append(stillWaitingGroups, g)
			}
		}

		if len(stillWaitingGroups) == 0 {
			delete(n.pendingGroups, key)
		} else {
			n.pendingGroups[key] = stillWaitingGroups
		}
	}
}

func (n *NotifyService) groupAlert(ctx context.Context, notifier config.NotifierSettings, a model.Alert) {
	_, span := otel.Tracer("").Start(ctx, "NotifyService.groupAlert")
	defer span.End()

	span.SetAttributes(attribute.String("alert", fmt.Sprintf("%+v", a)))

	key := map[string]string{}
	for _, l := range notifier.GroupLabels {
		key[l] = a.Labels[l]
	}

	n.groupMutex.Lock()
	groups, ok := n.pendingGroups[notifier.Name()]
	if !ok {
		groups = append(groups, groupMeta{
			GroupLabels: key,
			Timeout:     stubs.Time.Now().Add(notifier.GroupWait),
			Notifier:    notifier,
			Alerts:      []model.Alert{a},
		})
	} else {
		for i, g := range groups {
			if g.GroupLabels.Equal(key) {
				groups[i].Alerts = append(groups[i].Alerts, a)
				break
			}
		}
	}

	n.pendingGroups[notifier.Name()] = groups
	n.groupMutex.Unlock()
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
		if notifier.GroupWait != 0 {
			// If we have a GroupWait, we need to add this alert to a group.
			n.groupAlert(ctx, notifier, a)
			continue
		}

		if err := notifier.Notify(ctx, a); err != nil {
			n.bus.Logger("notify").Err(err).Msg("failed to notify for alert")
		}
	}

	if err := n.bus.Broadcaster().BroadcastAlerts(ctx, a); err != nil {
		n.bus.Logger("notify").Err(err).Msg("failed to broadcast the sucessful notify")
	}
}
