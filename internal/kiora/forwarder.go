package kiora

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type LocalForwarderProcessor struct{}

func (l *LocalForwarderProcessor) ProcessAlert(ctx context.Context, broadcast kioradb.ModelWriter, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "LocalForwarderProcessor.ProcessAlerts")
	defer span.End()

	return localdb.ProcessAlerts(ctx, *newAlert)
}

func (l *LocalForwarderProcessor) ProcessSilence(ctx context.Context, broadcast kioradb.ModelWriter, localdb kioradb.DB, silence *model.Silence) error {
	ctx, span := tracing.Tracer().Start(ctx, "LocalForwarderProcessor.ProcessSilence")
	defer span.End()

	return localdb.ProcessSilences(ctx, *silence)
}

type BroadcastProcessor struct{}

func (l *BroadcastProcessor) ProcessAlert(ctx context.Context, broadcast kioradb.ModelWriter, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "BroadcastProcessor.ProcessAlert")
	defer span.End()

	return broadcast.ProcessAlerts(ctx, *newAlert)
}

func (l *BroadcastProcessor) ProcessSilence(ctx context.Context, broadcast kioradb.ModelWriter, localdb kioradb.DB, silence *model.Silence) error {
	ctx, span := tracing.Tracer().Start(ctx, "BroadcastProcessor.ProcessSilence")
	defer span.End()

	return broadcast.ProcessSilences(ctx, *silence)
}
