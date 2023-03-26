package kioradb

import (
	"context"
	"sync"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ DB = &inMemoryDB{}

// inMemoryDB is a DB that does not persist anything, just storing all data in memory.
type inMemoryDB struct {
	aLock  sync.RWMutex
	alerts map[model.LabelsHash]model.Alert

	sLock    sync.RWMutex
	silences map[string]model.Silence
}

func NewInMemoryDB() *inMemoryDB {
	return &inMemoryDB{
		aLock:  sync.RWMutex{},
		alerts: make(map[model.LabelsHash]model.Alert),

		sLock:    sync.RWMutex{},
		silences: make(map[string]model.Silence),
	}
}

func (m *inMemoryDB) storeAlert(alert model.Alert) {
	labelsHash := alert.Labels.Hash()

	m.aLock.Lock()
	defer m.aLock.Unlock()
	m.alerts[labelsHash] = alert
}

func (m *inMemoryDB) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	for i := range alerts {
		m.storeAlert(alerts[i])
	}

	return nil
}

func (m *inMemoryDB) QueryAlerts(ctx context.Context, q *query.AlertQuery) []model.Alert {
	m.aLock.RLock()
	defer m.aLock.RUnlock()
	switch query := q.Filter.(type) {
	// Short circuit exact matches because we can process them more efficiently by just looking up the hash.
	case *query.ExactLabelMatchFilter:
		if existingAlert, ok := m.alerts[query.Labels.Hash()]; ok {
			return []model.Alert{existingAlert}
		}

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

func (m *inMemoryDB) StoreSilences(ctx context.Context, silences ...model.Silence) error {
	m.sLock.Lock()
	defer m.sLock.Unlock()
	for i := range silences {
		m.silences[silences[i].ID] = silences[i]
	}

	return nil
}

func (m *inMemoryDB) QuerySilences(ctx context.Context, query query.SilenceFilter) []model.Silence {
	m.sLock.RLock()
	defer m.sLock.RUnlock()
	silences := []model.Silence{}
	for _, silence := range m.silences {
		if query.MatchesSilence(ctx, &silence) {
			silences = append(silences, silence)
		}
	}

	return silences
}
