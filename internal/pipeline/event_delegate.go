package pipeline

import (
	"context"
	"time"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ = clustering.EventDelegate(&DBEventDelegate{})

// DBEventDelegate provides an EventDelegate that stores incoming data in the underlying database.
type DBEventDelegate struct {
	db kioradb.DB
}

func NewDBEventDelegate(db kioradb.DB) *DBEventDelegate {
	return &DBEventDelegate{
		db: db,
	}
}

func (d *DBEventDelegate) ProcessAlert(ctx context.Context, alert model.Alert) {
	currentAlerts := d.db.QueryAlerts(ctx, query.ExactLabelMatch(alert.Labels))

	// Copy attributes from the current alert if it exists.
	if len(currentAlerts) > 0 {
		currentAlert := currentAlerts[0]
		if alert.Status != model.AlertStatusResolved && alert.Status != model.AlertStatusTimedOut {
			if alert.LastNotifyTime.IsZero() {
				alert.LastNotifyTime = currentAlert.LastNotifyTime
			}
		}

		// If we have an alert coming from resolved or timed out back to firing, reset the last notify time so it'll notify again.
		if (currentAlert.Status == model.AlertStatusResolved || currentAlert.Status == model.AlertStatusTimedOut) && alert.Status == model.AlertStatusFiring {
			alert.LastNotifyTime = time.Time{}
		}

		if currentAlert.Acknowledgement != nil {
			alert.Acknowledgement = currentAlert.Acknowledgement
		}

		if currentAlert.Status == model.AlertStatusSilenced && alert.Status == model.AlertStatusFiring {
			alert.Status = model.AlertStatusSilenced
		}
	}

	// If it's firing, silence it if there's a matching silence. We can't do this async in a service
	// because that would cause a race condition where the alert could be fired before the silence is applied.
	if alert.Status == model.AlertStatusFiring {
		silences := d.db.QuerySilences(ctx, query.AllSilences(query.PartialLabelMatch(alert.Labels), query.SilenceIsActive()))
		if len(silences) > 0 {
			alert.Status = model.AlertStatusSilenced
		}
	}

	// TODO(cdouch): Handle errors here.
	d.db.StoreAlerts(ctx, alert) // nolint
}

func (d *DBEventDelegate) ProcessAlertAcknowledgement(ctx context.Context, alertID string, ack model.AlertAcknowledgement) {
	alerts := d.db.QueryAlerts(ctx, query.ID(alertID))
	if len(alerts) == 0 {
		// TODO(cdouch): Handle errors here.
		return
	}

	alert := alerts[0]
	alert.Acknowledgement = &ack

	if alert.Status == model.AlertStatusFiring {
		alert.Status = model.AlertStatusAcked
	}

	// TODO(cdouch): Handle errors here.
	d.db.StoreAlerts(ctx, alert) // nolint
}

func (d *DBEventDelegate) ProcessSilence(ctx context.Context, silence model.Silence) {
	existingSilence := d.db.QuerySilences(ctx, query.AllSilences(query.ID(silence.ID), query.SilenceIsActive()))
	if len(existingSilence) == 0 && silence.IsActive() {
		// This is a new silence, so we need to apply it to all the alerts.
		alerts := d.db.QueryAlerts(ctx, query.AlertFilterFunc(func(ctx context.Context, alert *model.Alert) bool {
			return silence.Matches(alert.Labels) && (alert.Status == model.AlertStatusFiring || alert.Status == model.AlertStatusAcked)
		}))

		for _, alert := range alerts {
			alert.Status = model.AlertStatusSilenced
			// TODO(cdouch): Handle errors here.
			d.db.StoreAlerts(ctx, alert) // nolint
		}
	}
	d.db.StoreSilences(ctx, silence) // nolint
}
