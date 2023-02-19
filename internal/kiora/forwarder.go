package kiora

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type LocalForwarderProcessor struct{}

func (l *LocalForwarderProcessor) ProcessAlert(ctx context.Context, broadcast kioradb.Broadcaster, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "LocalForwarderProcessor.ProcessAlerts")
	defer span.End()

	return localdb.StoreAlerts(ctx, *newAlert)
}

type BroadcastProcessor struct{}

func (l *BroadcastProcessor) ProcessAlert(ctx context.Context, broadcast kioradb.Broadcaster, localdb kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "BroadcastProcessor.ProcessAlert")
	defer span.End()

	return broadcast.BroadcastAlerts(ctx, *newAlert)
}
