package apiv1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/internal/server/api"
	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

var _ kioradb.DB = &mockDB{}

type mockDB struct {
	alerts   []model.Alert
	silences []model.Silence
}

func (m *mockDB) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	m.alerts = append(m.alerts, alerts...)
	return nil
}

func (m *mockDB) BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error {
	return m.StoreAlerts(ctx, alerts...)
}

func (m *mockDB) BroadcastAlertAcknowledgement(ctx context.Context, alertID string, ack model.AlertAcknowledgement) error {
	return nil
}

func (m *mockDB) QueryAlertStats(ctx context.Context, query query.AlertStatsQuery) ([]query.StatsResult, error) {
	return nil, nil
}

func (m *mockDB) BroadcastSilences(ctx context.Context, silences ...model.Silence) error {
	return nil
}

func (m *mockDB) QueryAlerts(ctx context.Context, query *query.AlertQuery) []model.Alert {
	return nil
}

func (m *mockDB) QuerySilences(ctx context.Context, query query.SilenceFilter) []model.Silence {
	return nil
}

func (m *mockDB) StoreSilences(ctx context.Context, silences ...model.Silence) error {
	m.silences = append(m.silences, silences...)
	return nil
}

func (m *mockDB) Close() error {
	return nil
}

func TestPostAlerts(t *testing.T) {
	// Construct a referenceTime that is used for each alert, and is expected to be found in the db.
	referenceTime, err := time.Parse(time.RFC3339, "2022-12-13T21:55:12Z")
	require.NoError(t, err)

	tests := []struct {
		name    string
		headers map[string]string
		body    []byte
	}{
		{
			name: "test json unmarshal",
			headers: map[string]string{
				"content-type": "application/json",
			},
			body: []byte(fmt.Sprintf(`[{
	"labels": {},
	"annotations": {},
	"status": "firing",
	"startsAt": "%s"
}]`, referenceTime.Format(time.RFC3339))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &mockDB{}
			api := New(
				api.NewAPIImpl(services.NewKioraBus(db, db, zerolog.New(os.Stderr), nil), nil),
				zerolog.New(os.Stderr),
			)

			request, err := http.NewRequest(http.MethodPost, "localhost/api/v1/alerts", bytes.NewReader(tt.body))
			require.NoError(t, err)

			for k, v := range tt.headers {
				request.Header.Add(k, v)
			}

			recorder := httptest.NewRecorder()

			api.postAlerts(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusAccepted, response.StatusCode, string(responseBody))

			require.Equal(t, 1, len(db.alerts), "expected one alert")
			alert := db.alerts[0]
			require.Equal(t, referenceTime, alert.StartTime)
			require.Equal(t, model.AlertStatusFiring, alert.Status)
		})
	}
}
