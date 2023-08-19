package regex

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
)

func init() {
	config.RegisterFilter("regex", NewRegexFilter)
}

// RegexFilter is a filter that matches if the given alert a) has the given label and b) that label matches a regex.
type RegexFilter struct {
	Label string         `config:"field" required:"true"`
	Regex *regexp.Regexp `config:"regex" required:"true"`
}

func NewRegexFilter(attrs map[string]string) (config.Filter, error) {
	delete(attrs, "type")
	var regexFilter RegexFilter
	if err := unmarshal.UnmarshalConfig(attrs, &regexFilter, unmarshal.UnmarshalOpts{DisallowUnknownFields: true}); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal regex filter")
	}

	return &regexFilter, nil
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
