package model_test

import (
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestMatchers(t *testing.T) {
	regexMatcher, err := model.LabelValueRegexMatcher("foo", "bar.+")
	require.NoError(t, err)
	testCases := []struct {
		name          string
		matcher       model.Matcher
		labels        model.Labels
		expectedMatch bool
	}{
		{
			name: "exact match",
			matcher: &model.LabelValueEqualMatcher{
				Label: "foo",
				Value: "bar",
			},
			labels: model.Labels{
				"foo": "bar",
			},
			expectedMatch: true,
		},
		{
			name: "exact match without label",
			matcher: &model.LabelValueEqualMatcher{
				Label: "foo",
				Value: "",
			},
			labels:        model.Labels{},
			expectedMatch: false,
		},
		{
			name:    "regex match",
			matcher: regexMatcher,
			labels: model.Labels{
				"foo": "barrington",
			},
			expectedMatch: true,
		},
		{
			name:    "negative match",
			matcher: model.NegativeMatcher(regexMatcher),
			labels: model.Labels{
				"foo": "barrington",
			},
			expectedMatch: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			match := tt.matcher.Match(tt.labels)
			if tt.expectedMatch && !match {
				t.Errorf("expected match, but it didn't")
			} else if !tt.expectedMatch && match {
				t.Errorf("expected a non match, but it did")
			}
		})
	}
}
