package apiv1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ kioradb.DB = &mockDB{}

type mockDB struct {
	alerts []model.Alert
}

func (m *mockDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	m.alerts = append(m.alerts, alerts...)
	return nil
}

func (m *mockDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return m.alerts, nil
}

func TestPostAlerts(t *testing.T) {
	// Construct a referenceTime that is used for each alert, and is expected to be found in the db.
	referenceTime, err := time.Parse(time.RFC3339, "2022-12-13T21:55:12Z")
	require.NoError(t, err)

	// Construct, then Marshal a kioraproto.Alert.
	msg := kioraproto.PostAlertsMessage{
		Alerts: []*kioraproto.Alert{
			{
				Status: kioraproto.AlertStatus_firing,
				StartTime: &timestamppb.Timestamp{
					Seconds: referenceTime.Unix(),
				},
			},
		},
	}

	alertBytes, err := proto.Marshal(&msg)
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
	"startTime": "%s"
}]`, referenceTime.Format(time.RFC3339))),
		},
		{
			name: "test proto unmarshal",
			headers: map[string]string{
				"content-type": "application/x-capnp",
			},
			body: alertBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &mockDB{}
			api := apiv1{db}

			request, err := http.NewRequest(http.MethodPost, "localhost/api/v1/alerts", bytes.NewReader(tt.body))
			require.NoError(t, err)

			for k, v := range tt.headers {
				request.Header.Add(k, v)
			}

			recorder := httptest.NewRecorder()

			api.postAlerts(recorder, request)

			response := recorder.Result()
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Equal(t, http.StatusAccepted, response.StatusCode, string(responseBody))

			require.Equal(t, 1, len(db.alerts), "expected one alert")
			alert := db.alerts[0]
			assert.Equal(t, referenceTime, alert.StartTime)
			assert.Equal(t, model.AlertStatusFiring, alert.Status)
		})
	}

}
