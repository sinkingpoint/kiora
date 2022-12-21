package model_test

import (
	"encoding/json"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func TestAlertUnmarshal(t *testing.T) {
	tests := []struct {
		name            string
		raw             string
		expectedFailure bool
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
				"status":"firing",
				"timeOutDeadline": "2022-12-21T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "start time after end time",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-22T21:32:27Z",
				"status":"firing",
				"timeOutDeadline": "2022-12-21T21:32:27Z"
}`,
			expectedFailure: true,
		},
		{
			name: "disallowed fields",
			raw: `{
				"labels":{},
				"annotations":{},
				"startTime": "2022-12-22T21:32:27Z",
				"status":"firing",
				"foo":"bar",
				"timeOutDeadline": "2022-12-21T21:32:27Z"
}`,
			expectedFailure: true,
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
		})
	}
}
