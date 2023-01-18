package kiora

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// SilenceApplier is an AlertProcessor that marks alerts as silenced if there is a silence that matches them.
type SilenceApplier struct {
}

func NewSilenceApplier() AlertProcessor {
	return &SilenceApplier{}
}

func (s *SilenceApplier) ProcessAlert(ctx context.Context, broadcast kioradb.ModelWriter, db kioradb.DB, existingAlert, newAlert *model.Alert) error {
	ctx, span := tracing.Tracer().Start(ctx, "SilenceApplier.ProcessAlert")
	defer span.End()

	if newAlert.Status == model.AlertStatusSilenced {
		span.AddEvent("Skipping because the alert is already silenced")
		return nil
	}

	if alreadySilenced := existingAlert != nil && existingAlert.Status == model.AlertStatusSilenced; alreadySilenced {
		span.AddEvent("Skipping because the existing alert is already silenced")
		newAlert.Status = model.AlertStatusSilenced
		return nil
	}

	silences, err := db.GetSilences(ctx, newAlert.Labels)
	if err != nil {
		return err
	}

	if len(silences) > 0 {
		newAlert.Status = model.AlertStatusSilenced

		if err := db.ProcessAlerts(ctx, *newAlert); err != nil {
			return err
		}
	}

	return nil
}
