package durationfilter_test

import (
	"context"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/durationfilter"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func TestDurationFilter(t *testing.T) {
	one := 1 * time.Second
	two := 5 * time.Hour
	tests := []struct {
		name            string
		filter          durationfilter.DurationFilter
		silence         model.Silence
		expectedSuccess bool
	}{
		{
			name: "test max duration",
			filter: durationfilter.DurationFilter{
				Field: "duration",
				Min:   nil,
				Max:   &one,
			},
			silence: model.Silence{
				StartTime: stubs.Time.Now(),
				EndTime:   stubs.Time.Now().Add(2 * time.Minute),
			},
			expectedSuccess: false,
		},
		{
			name: "test min duration",
			filter: durationfilter.DurationFilter{
				Field: "duration",
				Min:   &one,
			},
			silence: model.Silence{
				StartTime: stubs.Time.Now(),
				EndTime:   stubs.Time.Now().Add(2 * time.Minute),
			},
			expectedSuccess: true,
		},
		{
			name: "test both",
			filter: durationfilter.DurationFilter{
				Field: "duration",
				Min:   &one,
				Max:   &two,
			},
			silence: model.Silence{
				StartTime: stubs.Time.Now(),
				EndTime:   stubs.Time.Now().Add(2 * time.Minute),
			},
			expectedSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if success := tt.filter.Filter(context.Background(), &tt.silence); success != tt.expectedSuccess {
				t.Errorf("expected success: %t, but we didn't get it", tt.expectedSuccess)
			}
		})
	}
}
