package kiora

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// KioraProcessor is the main logic piece of Kiora that is responsible for actually acting on alerts, silences etc.
type KioraProcessor struct {
	DB              kioradb.DB
	Broadcaster     clustering.Broadcaster
	alertProcessors []AlertProcessor
}

// NewKioraProcessor creater a new KioraProcessor, starting the backing go routine that asynchronously processes incoming messages.
func NewKioraProcessor(db kioradb.DB, broadcaster clustering.Broadcaster) *KioraProcessor {
	processor := &KioraProcessor{
		DB:          db,
		Broadcaster: broadcaster,
	}

	return processor
}

func (k *KioraProcessor) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "KioraProcessor.ProcessAlerts")
	defer span.End()

	for _, alert := range alerts {
		if err := k.processAlert(ctx, alert); err != nil {
			log.Err(err).Msg("failed to process alert")
		}
	}

	return nil
}

func (k *KioraProcessor) QueryAlerts(ctx context.Context, query kioradb.AlertQuery) []model.Alert {
	return k.DB.QueryAlerts(ctx, query)
}

// AddAlertProcessor adds a processor to the stack of processors that get called when new alerts come in.
func (k *KioraProcessor) AddAlertProcessor(processor AlertProcessor) {
	k.alertProcessors = append(k.alertProcessors, processor)
}

func (k *KioraProcessor) processAlert(ctx context.Context, m model.Alert) error {
	var existingAlert *model.Alert
	if alerts := k.DB.QueryAlerts(ctx, &kioradb.ExactLabelMatchQuery{Labels: m.Labels}); len(alerts) > 0 {
		existingAlert = &alerts[0]
	}

	for _, processor := range k.alertProcessors {
		if err := processor.ProcessAlert(ctx, k.Broadcaster, k.DB, existingAlert, &m); err != nil {
			return err
		}
	}

	return nil
}

// AlertProcessor is a type that can be used to process an alert as it goes through the pipeline.
type AlertProcessor interface {
	ProcessAlert(ctx context.Context, broadcaster clustering.Broadcaster, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error
}

// AlertProcessorFunc wraps a func and turns it into an AlertProcessor.
type AlertProcessorFunc func(ctx context.Context, broadcaster clustering.Broadcaster, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error

func (a AlertProcessorFunc) ProcessAlert(ctx context.Context, broadcaster clustering.Broadcaster, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error {
	if a != nil {
		return a(ctx, broadcaster, localdb, existingAlert, newAlert)
	}

	return nil
}
