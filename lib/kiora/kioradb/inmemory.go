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

func (i *inMemoryDB) storeAlert(alert model.Alert) {
	labelsHash := alert.Labels.Hash()
	i.alerts[labelsHash] = alert
}

func (m *inMemoryDB) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	for i := range alerts {
		m.storeAlert(alerts[i])
	}

	return nil
}

func (m *inMemoryDB) QueryAlerts(ctx context.Context, query AlertQuery) []model.Alert {
	switch query := query.(type) {
	case *ExactLabelMatchQuery:
		if existingAlert, ok := m.alerts[query.Labels.Hash()]; ok {
			return []model.Alert{existingAlert}
		}

		// Short circuit exact matches because we can process them more efficiently by just looking up the hash.
		return []model.Alert{}
	default:
		alerts := []model.Alert{}
		for _, alert := range m.alerts {
			if query.MatchesAlert(ctx, &alert) {
				alerts = append(alerts, alert)
			}
		}
		return alerts
	}
}
