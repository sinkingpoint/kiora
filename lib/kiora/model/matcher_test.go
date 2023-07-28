package model_test

import (
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func MustRegexMatcher(t *testing.T, label, regex string) model.Matcher {
	t.Helper()
	m, err := model.LabelValueRegexMatcher(label, regex)
	require.NoError(t, err)
	return m
}

func ptrTo(m model.Matcher) *model.Matcher {
	return &m
}

func TestMatchers(t *testing.T) {
	testCases := []struct {
		name          string
		matcher       model.Matcher
		labels        model.Labels
		expectedMatch bool
	}{
		{
			name:    "exact match",
			matcher: model.LabelValueEqualMatcher("foo", "bar"),
			labels: model.Labels{
				"foo": "bar",
			},
			expectedMatch: true,
		},
		{
			name:          "exact match without label",
			matcher:       model.LabelValueEqualMatcher("foo", ""),
			labels:        model.Labels{},
			expectedMatch: false,
		},
		{
			name:    "regex match",
			matcher: MustRegexMatcher(t, "foo", "bar"),
			labels: model.Labels{
				"foo": "barrington",
			},
			expectedMatch: true,
		},
		{
			name:    "negative match",
			matcher: *ptrTo(MustRegexMatcher(t, "foo", "bar")).Negate(),
			labels: model.Labels{
				"foo": "barrington",
			},
			expectedMatch: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			match := tt.matcher.Matches(tt.labels)
			if tt.expectedMatch && !match {
				t.Errorf("expected match, but it didn't")
			} else if !tt.expectedMatch && match {
				t.Errorf("expected a non match, but it did")
			}
		})
	}
}
