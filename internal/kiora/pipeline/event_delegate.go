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
	alert.Status = model.AlertStatusAcked

	// TODO(cdouch): Handle errors here.
	d.db.StoreAlerts(ctx, alert) // nolint
}
