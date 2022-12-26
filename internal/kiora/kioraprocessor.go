package kiora

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ kioradb.DB = &KioraProcessor{}

// ErrProcessorClosed it the error returned when new alerts come in but the processor has been shut down.
var ErrProcessorClosed = errors.New("kioraprocessor is closed")

// AlertProcessor is a type that can be used to process an alert as it goes through the pipeline.
type AlertProcessor interface {
	Exec(ctx context.Context, db kioradb.DB, existingAlert, newAlert *model.Alert) error
}

// SilenceProcessor is a type that can be used to process a silence as it goes through the pipeline.
type SilenceProcessor interface {
	Exec(ctx context.Context, db kioradb.DB, silence *model.Silence)
}

// KioraProcessor is the main logic piece of Kiora that is responsible for actually acting on alerts, silences etc.
type KioraProcessor struct {
	alertProcessors   []AlertProcessor
	silenceProcessors []SilenceProcessor
	backingDB         kioradb.DB

	killChannel        chan struct{}
	killed             bool
	processingPipeline chan any
}

// NewKioraProcessor creater a new KioraProcessor, starting the backing go routine that asynchronously processes incoming messages.
func NewKioraProcessor(db kioradb.DB) *KioraProcessor {
	processor := &KioraProcessor{
		killChannel:        make(chan struct{}),
		killed:             false,
		processingPipeline: make(chan any, 100), // TODO(cdouch): make the queue length configurable.
		backingDB:          db,
	}

	go func() {
	outer:
		for {
			select {
			case <-processor.killChannel:
				break outer
			case m := <-processor.processingPipeline:
				processor.process(m)
			}
		}
	}()

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

// Kill closes the Processor, which will cause any incoming alerts to fail.
func (k *KioraProcessor) Kill() {
	k.killed = true
	close(k.killChannel)
}

func (k *KioraProcessor) process(m any) {
	switch m := m.(type) {
	case []model.Alert:
		for _, m := range m {
			k.processAlert(m)
		}
	case []model.Silence:
		for _, m := range m {
			k.processSilence(m)
		}
	default:
		panic(fmt.Sprintf("BUG: unhandled type in the processing pipeline: %T", m))
	}
}

func (k *KioraProcessor) processAlert(m model.Alert) {
	ctx := context.Background()
	existingAlert, err := k.GetExistingAlert(ctx, m.Labels)
	if err != nil {
		log.Err(err).Msg("failed to get existing alerts from backend. Dropping alert.")
		return
	}

	for _, processor := range k.alertProcessors {
		processor.Exec(ctx, k.backingDB, existingAlert, &m)
	}
}

func (k *KioraProcessor) processSilence(m model.Silence) {
	ctx := context.Background()

	for _, processor := range k.silenceProcessors {
		processor.Exec(ctx, k.backingDB, &m)
	}
}

func (k *KioraProcessor) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	if k.killed {
		return ErrProcessorClosed
	}
	k.processingPipeline <- alerts
	return k.backingDB.ProcessAlerts(ctx, alerts...)
}

func (k *KioraProcessor) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	if k.killed {
		return ErrProcessorClosed
	}
	k.processingPipeline <- silences
	return k.backingDB.ProcessSilences(ctx, silences...)
}

func (k *KioraProcessor) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return k.backingDB.GetAlerts(ctx)
}

func (k *KioraProcessor) GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error) {
	return k.backingDB.GetExistingAlert(ctx, labels)
}
