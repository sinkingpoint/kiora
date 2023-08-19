package regex

import (
	"context"

	"github.com/grafana/regexp"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
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
	if err != nil {
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

func (r *RegexFilter) Describe() string {
	return "field " + r.Label + " doesn't match " + r.Regex.String()
}

func (r *RegexFilter) Filter(ctx context.Context, f config.Fielder) bool {
	value, err := f.Field(r.Label)
	if err != nil {
		return false
	}

	if label, ok := value.(string); ok {
		return r.Regex.MatchString(label)
	}

	log.Warn().Str("field", r.Label).Interface("value", value).Msg("regex filter: field is not a string")
	return false
}