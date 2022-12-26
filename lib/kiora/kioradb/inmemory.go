package kioradb

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ DB = &inMemoryDB{}

type AlertIngestor = func(db DB, existingAlert, newAlert *model.Alert)

// inMemoryDB is a DB that does not persist anything, just storing all data in memory.
type inMemoryDB struct {
	alerts   map[model.LabelsHash]model.Alert
	silences map[string]model.Silence
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

func (m *inMemoryDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	for i := range silences {
		silence := &silences[i]
		m.silences[silence.ID] = *silence
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

func (m *inMemoryDB) GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error) {
	hash := labels.Hash()
	if existingAlert, ok := m.alerts[hash]; ok {
		return &existingAlert, nil
	}

	return nil, nil
}

func (m *inMemoryDB) GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error) {
	silences := []model.Silence{}
	for _, silence := range m.silences {
		if silence.Matches(labels) {
			silences = append(silences, silence)
		}
	}

	return silences, nil
}
