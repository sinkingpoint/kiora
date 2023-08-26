package regex

import (
	"context"
	"fmt"

	"github.com/grafana/regexp"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
)

// RegexFilter is a filter that matches if the given alert a) has the given label and b) that label matches a regex.
type RegexFilter struct {
	Label string         `config:"field" required:"true"`
	Regex *regexp.Regexp `config:"regex" required:"true"`
}

func NewFilter(globals *config.Globals, attrs map[string]string) (config.Filter, error) {
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

func (r *RegexFilter) Filter(ctx context.Context, f config.Fielder) error {
	value, err := f.Field(r.Label)
	if err != nil {
		return fmt.Errorf("failed to get field %q: %w", r.Label, err)
	}

	if label, ok := value.(string); ok {
		if r.Regex.MatchString(label) {
			return nil
		}

		return fmt.Errorf("label %q does not match regex %q", label, r.Regex.String())
	}

	return fmt.Errorf("label %q is not a string", r.Label)
}
