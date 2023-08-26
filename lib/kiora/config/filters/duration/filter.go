package duration

import (
	"context"
	"fmt"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
)

type DurationFilter struct {
	Field string         `config:"field" required:"true"`
	Max   *time.Duration `config:"max"`
	Min   *time.Duration `config:"min"`
}

func NewFilter(globals *config.Globals, attrs map[string]string) (config.Filter, error) {
	delete(attrs, "type")
	var durationFilter DurationFilter

	if err := unmarshal.UnmarshalConfig(attrs, &durationFilter, unmarshal.UnmarshalOpts{DisallowUnknownFields: true}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal duration filter: %w", err)
	}

	if durationFilter.Max == nil && durationFilter.Min == nil {
		return nil, fmt.Errorf("duration filter must have at least one of max or min")
	}

	return &durationFilter, nil
}

func (d *DurationFilter) Type() string {
	return "duration"
}

func (d *DurationFilter) Filter(ctx context.Context, fielder config.Fielder) error {
	field, err := fielder.Field(d.Field)
	if err != nil {
		return fmt.Errorf("failed to get field %q: %w", d.Field, err)
	}

	duration, ok := field.(time.Duration)
	if !ok {
		return fmt.Errorf("field %q is not a duration", d.Field)
	}

	if d.Max != nil && duration > *d.Max {
		return fmt.Errorf("field %q is greater than %s", d.Field, d.Max.String())
	}

	if d.Min != nil && duration < *d.Min {
		return fmt.Errorf("field %q is less than %s", d.Field, d.Min.String())
	}

	return nil
}
