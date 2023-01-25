package kiora

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel/attribute"
)

// NotifierConfig is an interface that defines a Configuration for the NotifierProcessor.
type NotifierConfig interface {
	// GetNotifiersForAlert takes an Alert and returns all the notifiers that should be notified about the alert.
	GetNotifiersForAlert(a *model.Alert) []kioradb.ModelWriter
}

// NotifierProcessor is an Alert Processor responsible for actually notifying for alerts.
type NotifierProcessor struct {
	me     string
	config NotifierConfig
}

func NewNotifierProcessor(myName string, config NotifierConfig) *NotifierProcessor {
	return &NotifierProcessor{
		me:     myName,
		config: config,
	}
}

func (n *NotifierProcessor) ProcessAlert(ctx context.Context, broadcast kioradb.ModelWriter, db kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "NotifierProcessor.ProcessAlert")
	defer span.End()

	span.SetAttributes(attribute.String("alert", fmt.Sprintf("%+v", newAlert)))

	// Skip this notify if we're not authoritative in the new alert, or the existing alert.
	if newAlert.AuthNode != n.me || (existingAlert != nil && existingAlert.AuthNode != n.me) {
		span.AddEvent("Skipping because this node is not authoritative")
		return nil
	}

	if newAlert.Status != model.AlertStatusProcessing || (existingAlert != nil && existingAlert.Status != model.AlertStatusProcessing) {
		span.AddEvent("Skipping because the alert isn't processing")
		return nil
	}

	newAlert.Status = model.AlertStatusFiring
	var notifyError error
	notifiers := n.config.GetNotifiersForAlert(newAlert)
	for _, notify := range notifiers {
		if err := notify.ProcessAlerts(ctx, *newAlert); err != nil {
			notifyError = multierror.Append(notifyError, err)
		}
	}

	if err := db.ProcessAlerts(ctx, *newAlert); err != nil {
		notifyError = multierror.Append(notifyError, err)
	}

	if err := broadcast.ProcessAlerts(ctx, *newAlert); err != nil {
		notifyError = multierror.Append(notifyError, err)
	}

	return notifyError
}
