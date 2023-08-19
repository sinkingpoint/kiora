package duration

import (
	"context"
	"fmt"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
)

func init() {
	config.RegisterFilter("duration", NewDurationFilter)
}

type DurationFilter struct {
	Field string         `config:"field" required:"true"`
	Max   *time.Duration `config:"max"`
	Min   *time.Duration `config:"min"`
}

func NewDurationFilter(attrs map[string]string) (config.Filter, error) {
	delete(attrs, "type")
	var durationFilter DurationFilter

	if err := unmarshal.UnmarshalConfig(attrs, &durationFilter, unmarshal.UnmarshalOpts{DisallowUnknownFields: true}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal duration filter: %w", err)
	}

	return &durationFilter, nil
}

func (d *DurationFilter) Describe() string {
	if d.Max != nil && d.Min != nil {
		return fmt.Sprintf("field %q must be between %s and %s", d.Field, d.Min.String(), d.Max.String())
	}

	if d.Max != nil {
		return fmt.Sprintf("field %q must be less than %s", d.Field, d.Max.String())
	}

	if d.Min != nil {
		return fmt.Sprintf("field %q must be greater than %s", d.Field, d.Min.String())
	}

	panic("BUG: DurationFilter has neither a max nor a min")
}

func (d *DurationFilter) Type() string {
	return "duration"
}

func (d *DurationFilter) Filter(ctx context.Context, fielder config.Fielder) bool {
	field, err := fielder.Field(d.Field)
	if err != nil {
		return false
	}

	duration, ok := field.(time.Duration)
	if !ok {
		return false
	}

	if d.Max != nil && duration > *d.Max {
		return false
	}

	if d.Min != nil && duration < *d.Min {
		return false
	}

	return true
}
