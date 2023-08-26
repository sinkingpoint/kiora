package regex_test

import (
	"context"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/regex"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestRegexFilter(t *testing.T) {
	tests := []struct {
		Name        string
		Label       string
		Regex       string
		Alert       model.Alert
		ShouldMatch bool
	}{
		{
			Name:  "standard match",
			Label: "test",
			Regex: "test",
			Alert: model.Alert{
				Labels: model.Labels{
					"test": "test",
				},
			},
			ShouldMatch: true,
		},
		{
			Name:  "non existent label",
			Label: "some_weird_non_existent_label",
			Regex: "test",
			Alert: model.Alert{
				Labels: model.Labels{
					"test": "test",
				},
			},
			ShouldMatch: false,
		},
		{
			Name:  "non regex match",
			Label: "test",
			Regex: "^test$",
			Alert: model.Alert{
				Labels: model.Labels{
					"test": "not test",
				},
			},
			ShouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			filter, err := regex.NewFilter(nil, map[string]string{
				"field": tt.Label,
				"regex": tt.Regex,
			})

			require.NoError(t, err)

			err = filter.Filter(context.TODO(), &tt.Alert)
			if tt.ShouldMatch {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
