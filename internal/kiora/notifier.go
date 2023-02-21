package kiora

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/kiora/lib/kiora/notify"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// NotifierConfig is an interface that defines a Configuration for the NotifierProcessor.
type NotifierConfig interface {
	// GetNotifiersForAlert takes an Alert and returns all the notifiers that should be notified about the alert.
	GetNotifiersForAlert(a *model.Alert) []notify.Notifier
}

// NotifierProcessor is an Alert Processor responsible for actually notifying for alerts.
type NotifierProcessor struct {
	config NotifierConfig
}

func NewNotifierProcessor(config NotifierConfig) *NotifierProcessor {
	return &NotifierProcessor{
		config: config,
	}
}

func (n *NotifierProcessor) ProcessAlert(ctx context.Context, broadcaster kioradb.Broadcaster, db kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "NotifierProcessor.ProcessAlert")
	defer span.End()

	span.SetAttributes(attribute.String("alert", fmt.Sprintf("%+v", newAlert)))

	// Before we send any notifications, if this is a new alert, or an update to the states, save it in the local db.
	if existingAlert == nil || newAlert.Status != model.AlertStatusProcessing {
		if err := db.StoreAlerts(ctx, *newAlert); err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}

	if newAlert.Status != model.AlertStatusProcessing || (existingAlert != nil && existingAlert.Status != model.AlertStatusProcessing) {
		span.AddEvent("Skipping because the alert isn't processing")
		return nil
	}

	newAlert.Status = model.AlertStatusFiring
	var notifyError error
	notifiers := n.config.GetNotifiersForAlert(newAlert)
	for _, notify := range notifiers {
		if err := notify.Notify(ctx, *newAlert); err != nil {
			notifyError = multierror.Append(notifyError, err)
		}
	}

	if err := db.StoreAlerts(ctx, *newAlert); err != nil {
		notifyError = multierror.Append(notifyError, err)
	}

	if err := broadcaster.BroadcastAlerts(ctx, *newAlert); err != nil {
		notifyError = multierror.Append(notifyError, err)
	}

	return notifyError
}
