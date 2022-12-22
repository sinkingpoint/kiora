package model_test

import (
	"encoding/json"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func TestSilenceUnmarshal(t *testing.T) {
	tests := []struct {
		name            string
		raw             string
		expectedFailure bool
	}{
		{
			name: "standard",
			raw: `{
				"startTime": "2022-12-21T21:32:27Z",
				"endTime": "2022-12-22T21:32:27Z",
				"matchers": [
					{
						"label": "foo",
						"value": "bar",
						"negative": false,
						"regex": false
					}
				]
}`,
			expectedFailure: false,
		},
		{
			name: "empty matchers",
			raw: `{
				"startTime": "2022-12-21T21:32:27Z",
				"endTime": "2022-12-22T21:32:27Z",
				"matchers": []
}`,
			expectedFailure: true,
		},
		{
			name: "missing start time",
			raw: `{
				"endTime": "2022-12-22T21:32:27Z",
				"matchers": [{"label": "foo","value": "bar","negative": false,"regex": false}]
}`,
			expectedFailure: true,
		},
		{
			name: "missing end time",
			raw: `{
				"startTime": "2022-12-21T21:32:27Z",
				"matchers": [{"label": "foo","value": "bar","negative": false,"regex": false}]
}`,
			expectedFailure: true,
		},
		{
			name: "end time not after start time",
			raw: `{
				"startTime": "2022-12-22T21:32:27Z",
				"endTime": "2022-12-21T21:32:27Z",
				"matchers": [{"label": "foo","value": "bar","negative": false,"regex": false}]
}`,
			expectedFailure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.raw)
			var alert model.Silence

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
