package config

import (
	"errors"
	"regexp"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type Filter interface {
	Type() string
	FilterAlert(a *model.Alert) bool
}

type Link struct {
	incomingFilter Filter
	to             string
}

type FilterConstructor = func(n edge) (Filter, error)

var filterRegistry = map[string]FilterConstructor{
	"regex": NewRegexFilter,
}

func LookupFilter(name string) FilterConstructor {
	return filterRegistry[name]
}

// RegexFilter is a filter that matches if the given alert a) has the given label and b) that label matches a regex.
type RegexFilter struct {
	Label string
	Regex *regexp.Regexp
}

func NewRegexFilter(n edge) (Filter, error) {
	field := n.attrs["field"]
	if field == "" {
		return nil, errors.New("expected `field` in regex filter")
	}

	regexStr := n.attrs["regex"]
	if regexStr == "" {
		return nil, errors.New("expected `regex` in regex filter")
	}

	regex, err := regexp.Compile(regexStr)
	if err == nil {
		return nil, err
	}

	return &RegexFilter{
		Label: field,
		Regex: regex,
	}, nil
}

func (r *RegexFilter) Type() string {
	return "regex"
}

func (r *RegexFilter) FilterAlert(a *model.Alert) bool {
	if label, ok := a.Labels[r.Label]; ok {
		return r.Regex.MatchString(label)
	}

	return false
}
