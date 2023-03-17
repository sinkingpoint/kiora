package regexfilter

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func init() {
	config.RegisterFilter("regex", NewRegexFilter)
}

// RegexFilter is a filter that matches if the given alert a) has the given label and b) that label matches a regex.
type RegexFilter struct {
	Label string
	Regex *regexp.Regexp
}

func NewRegexFilter(attrs map[string]string) (config.Filter, error) {
	field := attrs["field"]
	if field == "" {
		return nil, errors.New("expected `field` in regex filter")
	}

	regexStr := attrs["regex"]
	if regexStr == "" {
		return nil, errors.New("expected `regex` in regex filter")
	}

	regex, err := regexp.Compile(regexStr)
	if err == nil {
		return nil, errors.Wrap(err, "failed to compile regex")
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
