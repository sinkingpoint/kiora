package kiora

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ kioradb.DB = &KioraProcessor{}

// AlertProcessor is a type that can be used to process an alert as it goes through the pipeline.
type AlertProcessor interface {
	ProcessAlert(ctx context.Context, broadcast kioradb.ModelWriter, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error
}

// SilenceProcessor is a type that can be used to process a silence as it goes through the pipeline.
type SilenceProcessor interface {
	ProcessSilence(ctx context.Context, broadcast kioradb.ModelWriter, db kioradb.DB, silence *model.Silence) error
}

// KioraProcessor is the main logic piece of Kiora that is responsible for actually acting on alerts, silences etc.
type KioraProcessor struct {
	*kioradb.FallthroughDB
	Broadcast         kioradb.ModelWriter
	alertProcessors   []AlertProcessor
	silenceProcessors []SilenceProcessor
}

// NewKioraProcessor creater a new KioraProcessor, starting the backing go routine that asynchronously processes incoming messages.
func NewKioraProcessor(db kioradb.DB, broadcaster kioradb.ModelWriter) *KioraProcessor {
	processor := &KioraProcessor{
		FallthroughDB: kioradb.NewFallthroughDB(db),
		Broadcast:     broadcaster,
	}

	return processor
}

// AddAlertProcessor adds a processor to the stack of processors that get called when new alerts come in.
func (k *KioraProcessor) AddAlertProcessor(processor AlertProcessor) {
	k.alertProcessors = append(k.alertProcessors, processor)
}

// AddSilenceProccessor adds a processor to the stack of processors that get called when new silences come in.
func (k *KioraProcessor) AddSilenceProccessor(processor SilenceProcessor) {
	k.silenceProcessors = append(k.silenceProcessors, processor)
}

func (k *KioraProcessor) processAlert(ctx context.Context, m model.Alert) error {
	existingAlert, err := k.GetExistingAlert(ctx, m.Labels)
	if err != nil {
		return err
	}

	for _, processor := range k.alertProcessors {
		if err := processor.ProcessAlert(ctx, k.Broadcast, k.FallthroughDB, existingAlert, &m); err != nil {
			return err
		}
	}

	return nil
}

func (k *KioraProcessor) processSilence(ctx context.Context, m model.Silence) error {
	for _, processor := range k.silenceProcessors {
		if err := processor.ProcessSilence(ctx, k.Broadcast, k.FallthroughDB, &m); err != nil {
			return err
		}
	}

	return nil
}

func (k *KioraProcessor) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "KioraProcessor.ProcessAlerts")
	defer span.End()

	for _, alert := range alerts {
		if err := k.processAlert(ctx, alert); err != nil {
			log.Err(err).Msg("failed to process alert")
		}
	}

	return nil
}

func (k *KioraProcessor) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	ctx, span := tracing.Tracer().Start(ctx, "KioraProcessor.ProcessSilences")
	defer span.End()

	for _, silence := range silences {
		if err := k.processSilence(ctx, silence); err != nil {
			log.Err(err).Msg("failed to process silence")
		}
	}

	return nil
}
