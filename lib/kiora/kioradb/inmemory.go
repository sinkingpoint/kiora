package kioradb

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ DB = &inMemoryDB{}

// inMemoryDB is a DB that does not persist anything, just storing all data in memory.
type inMemoryDB struct {
	alerts map[model.LabelsHash]model.Alert
}

func NewInMemoryDB() *inMemoryDB {
	return &inMemoryDB{
		alerts: make(map[model.LabelsHash]model.Alert),
	}
}

func (i *inMemoryDB) processAlert(alert model.Alert) {
	labelsHash := alert.Labels.Hash()

	// The only thing we change on incoming alerts is to extend the start time back
	// if we have one before its start time.
	if existingAlert, hasAlert := i.alerts[labelsHash]; hasAlert {
		if existingAlert.StartTime.Before(alert.StartTime) {
			alert.StartTime = existingAlert.StartTime
		}
	}

	i.alerts[labelsHash] = alert
}

func (m *inMemoryDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	for i := range alerts {
		m.processAlert(alerts[i])
	}

	return nil
}
