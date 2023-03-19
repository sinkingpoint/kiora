package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertUnmarshal(t *testing.T) {
	referenceStartTime, err := time.Parse(time.RFC3339, "2022-12-21T21:32:27Z")
	require.NoError(t, err)

	referenceTimeoutDeadline, err := time.Parse(time.RFC3339, "2022-12-22T21:32:27Z")
	require.NoError(t, err)

	referenceNowTime := time.Now()
	stubs.Time.Now = func() time.Time {
		return referenceNowTime
	}

	tests := []struct {
		name            string
		raw             string
		expectedFailure bool
		expectedAlert   *model.Alert
	}{
		{
			name: "standard",
			raw: `{
				"labels":{
					"foo": "bar"
				},
				"annotations":{
					"bar": "baz"
				},
				"startTime": "2022-12-21T21:32:27Z",
				"status":"firing",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: false,
			expectedAlert: &model.Alert{
				Labels: model.Labels(map[string]string{
					"foo": "bar",
				}),
				Annotations: map[string]string{
					"bar": "baz",
				},
				StartTime:       referenceStartTime,
				Status:          model.AlertStatusFiring,
				TimeOutDeadline: referenceTimeoutDeadline,
			},
		},
		{
			name: "missing labels",
			raw: `{
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"status":"firing",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "missing annotations",
			raw: `{
				"labels":{},
				"startTime": "2022-12-21T21:32:27Z",
				"status":"firing",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "missing status",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "missing startTime",
			raw: `{
				"labels":{},
				"annotations":{},
				"status":"firing",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "start time == end time",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"endTime": "2022-12-21T21:32:27Z",
				"status":"firing",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "start time after end time",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"endTime": "2022-12-20T21:32:27Z",
				"status":"firing",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "disallowed fields",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"status":"firing",
				"foo":"bar",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "start time unset",
			raw: `{
				"labels":{},
				"annotations":{},
				"status":"firing"
			}`,
			expectedFailure: false,
			expectedAlert: &model.Alert{
				Labels:          make(model.Labels),
				Annotations:     make(map[string]string),
				StartTime:       referenceNowTime,
				Status:          model.AlertStatusFiring,
				TimeOutDeadline: referenceNowTime.Add(model.DEFAULT_TIMEOUT_INTERVAL),
			},
		},
		{
			name: "resolved but unset endtime",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"status":"resolved",
				"timeOutDeadline": "2022-12-22T21:32:27Z"
			}`,
			expectedFailure: false,
			expectedAlert: &model.Alert{
				Labels:          make(model.Labels),
				Annotations:     make(map[string]string),
				StartTime:       referenceStartTime,
				Status:          model.AlertStatusResolved,
				TimeOutDeadline: referenceTimeoutDeadline,
				EndTime:         referenceNowTime,
			},
		},
		{
			name: "timeout deadline unset",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-21T21:32:27Z",
				"status":"firing"
			}`,
			expectedFailure: false,
			expectedAlert: &model.Alert{
				Labels:          make(model.Labels),
				Annotations:     make(map[string]string),
				StartTime:       referenceStartTime,
				Status:          model.AlertStatusFiring,
				TimeOutDeadline: referenceStartTime.Add(model.DEFAULT_TIMEOUT_INTERVAL),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.raw)
			var alert model.Alert

			err := json.Unmarshal(raw, &alert)
			isError := err != nil

			if !isError && tt.expectedFailure {
				t.Errorf("expected an error, but didn't get one")
			} else if isError && !tt.expectedFailure {
				t.Errorf("didn't expect an error, but got %q", err.Error())
			}

			if !tt.expectedFailure && tt.expectedAlert != nil {
				assert.Equal(t, *tt.expectedAlert, alert)
			}
		})
	}
}
