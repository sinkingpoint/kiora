package kiora

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// SilenceApplier is an AlertProcessor that marks alerts as silenced if there is a silence that matches them.
type SilenceApplier struct {
}

func NewSilenceApplier() AlertProcessor {
	return &SilenceApplier{}
}

func (s *SilenceApplier) Exec(ctx context.Context, db kioradb.DB, existingAlert, newAlert *model.Alert) error {
	if alreadySilenced := existingAlert != nil && existingAlert.Status == model.AlertStatusSilenced; alreadySilenced {
		return nil
	}

	silences, err := db.GetSilences(ctx, newAlert.Labels)
	if err != nil {
		return err
	}

	if len(silences) > 0 {
		newAlert.Status = model.AlertStatusSilenced
	}

	return db.ProcessAlerts(ctx, *newAlert)
}
