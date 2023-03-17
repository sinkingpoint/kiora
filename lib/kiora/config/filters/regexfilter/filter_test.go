package regexfilter_test

import (
	"regexp"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/regexfilter"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
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
			Name:  "non existant label",
			Label: "some_weird_non_existant_label",
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
			regex := regexp.MustCompile(tt.Regex)
			match := regexfilter.RegexFilter{
				Label: tt.Label,
				Regex: regex,
			}

			matches := match.FilterAlert(&tt.Alert)
			assert.Equal(t, tt.ShouldMatch, matches)
		})
	}
}
