package kioradb

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
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

func (m *inMemoryDB) GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error) {
	silences := []model.Silence{}
	for _, silence := range m.silences {
		if silence.Matches(labels) {
			silences = append(silences, silence)
		}
	}

	return silences, nil
}

func (f *inMemoryDB) RegisterEndpoints(ctx context.Context, httpRouter *mux.Router, grpcServer *grpc.Server) error {
	return nil
}
