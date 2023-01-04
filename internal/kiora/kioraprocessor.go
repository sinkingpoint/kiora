package kiora

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ kioradb.DB = &KioraProcessor{}

// AlertProcessor is a type that can be used to process an alert as it goes through the pipeline.
type AlertProcessor interface {
	Exec(ctx context.Context, db kioradb.DB, existingAlert, newAlert *model.Alert) error
}

// SilenceProcessor is a type that can be used to process a silence as it goes through the pipeline.
type SilenceProcessor interface {
	Exec(ctx context.Context, db kioradb.DB, silence *model.Silence) error
}

// KioraProcessor is the main logic piece of Kiora that is responsible for actually acting on alerts, silences etc.
type KioraProcessor struct {
	*kioradb.FallthroughDB
	alertProcessors   []AlertProcessor
	silenceProcessors []SilenceProcessor
}

// NewKioraProcessor creater a new KioraProcessor, starting the backing go routine that asynchronously processes incoming messages.
func NewKioraProcessor(db kioradb.DB) *KioraProcessor {
	processor := &KioraProcessor{
		FallthroughDB: kioradb.NewFallthroughDB(db),
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
		if err := processor.Exec(ctx, k.FallthroughDB, existingAlert, &m); err != nil {
			return err
		}
	}

	return nil
}

func (k *KioraProcessor) processSilence(ctx context.Context, m model.Silence) error {
	for _, processor := range k.silenceProcessors {
		if err := processor.Exec(ctx, k.FallthroughDB, &m); err != nil {
			return err
		}
	}

	return nil
}

func (k *KioraProcessor) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	for _, alert := range alerts {
		if err := k.processAlert(ctx, alert); err != nil {
			log.Err(err).Msg("failed to process alert")
		}
	}

	return nil
}

func (k *KioraProcessor) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	for _, silence := range silences {
		if err := k.processSilence(ctx, silence); err != nil {
			log.Err(err).Msg("failed to process silence")
		}
	}

	return nil
}
