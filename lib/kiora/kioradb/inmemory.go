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
	i.alerts[labelsHash] = alert
}

func (m *inMemoryDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	for i := range alerts {
		m.processAlert(alerts[i])
	}

	return nil
}

func (m *inMemoryDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	alerts := make([]model.Alert, 0, len(m.alerts))
	for _, v := range m.alerts {
		alerts = append(alerts, v)
	}
	return alerts, nil
}
