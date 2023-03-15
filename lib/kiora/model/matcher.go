package model

import (
	"regexp"
)

var _ Matcher = &LabelValueEqualMatcher{}
var _ Matcher = &labelValueRegexMatcher{}
var _ Matcher = &negativeMatcher{}

// Matcher is an interface that defines something that can be used to match, or reject a labelset.
type Matcher interface {
	Match(labels Labels) bool
}

type LabelValueEqualMatcher struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (l *LabelValueEqualMatcher) Match(labels Labels) bool {
	if label, ok := labels[l.Label]; ok {
		return label == l.Value
	}

	// Never match if the label doesn't exist. This is a departure from Alertmanager,
	// where a non existant label is treated as an empty string, and causes footguns when we encouter a silence
	// that silences everything without a label (weird_label="")
	return false
}

type labelValueRegexMatcher struct {
	Label string `json:"label"`
	Value string `json:"value"`
	regex *regexp.Regexp
}

func LabelValueRegexMatcher(label string, value string) (*labelValueRegexMatcher, error) {
	regex, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}

	return &labelValueRegexMatcher{
		Label: label,
		Value: value,
		regex: regex,
	}, nil
}

func (l *labelValueRegexMatcher) Match(labels Labels) bool {
	if label, ok := labels[l.Label]; ok {
		return l.regex.MatchString(label)
	}

	return false
}

type negativeMatcher struct {
	matcher Matcher
}

func NegativeMatcher(matcher Matcher) Matcher {
	return &negativeMatcher{
		matcher: matcher,
	}
}

func (l *negativeMatcher) Match(labels Labels) bool {
	return !l.matcher.Match(labels)
}
