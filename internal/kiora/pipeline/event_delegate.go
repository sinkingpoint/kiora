package pipeline

import (
	"context"

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
	currentAlert := d.db.QueryAlerts(ctx, &query.ExactLabelMatchQuery{Labels: alert.Labels})
	if len(currentAlert) > 0 {
		// If we have the alert already, and this new alert is in Processing (i.e. it just came in), skip it.
		// This works to deduplicate alerts.
		if alert.Status == model.AlertStatusProcessing {
			return
		}
	}

	// TODO(cdouch): Handle errors here.
	d.db.StoreAlerts(ctx, alert) // nolint
}
