package durationfilter

import (
	"context"
	"fmt"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

func init() {
	config.RegisterFilter("duration", NewDurationFilter)
}

type DurationFilter struct {
	Field string
	Max   *time.Duration
	Min   *time.Duration
}

func NewDurationFilter(attrs map[string]string) (config.Filter, error) {
	field := attrs["field"]
	if field == "" {
		return nil, fmt.Errorf("expected `field` in duration filter")
	}

	durations := &DurationFilter{
		Field: field,
	}

	minStr := attrs["min"]
	if minStr != "" {
		minDuration, err := time.ParseDuration(minStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse min duration: %w", err)
		}

		durations.Min = &minDuration
	}

	maxStr := attrs["max"]
	if maxStr != "" {
		maxDuration, err := time.ParseDuration(maxStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse max duration: %w", err)
		}

		durations.Max = &maxDuration
	}

	return durations, nil
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
