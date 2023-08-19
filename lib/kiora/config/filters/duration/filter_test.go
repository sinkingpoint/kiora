package duration_test

import (
	"context"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/duration"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestDurationFilter(t *testing.T) {
	tests := []struct {
		name            string
		attrs           map[string]string
		silence         model.Silence
		expectedSuccess bool
	}{
		{
			name: "test max duration",
			attrs: map[string]string{
				"field": "duration",
				"max":   "1s",
			},
			silence: model.Silence{
				StartTime: stubs.Time.Now(),
				EndTime:   stubs.Time.Now().Add(2 * time.Minute),
			},
			expectedSuccess: false,
		},
		{
			name: "test min duration",
			attrs: map[string]string{
				"field": "duration",
				"min":   "1s",
			},
			silence: model.Silence{
				StartTime: stubs.Time.Now(),
				EndTime:   stubs.Time.Now().Add(2 * time.Minute),
			},
			expectedSuccess: true,
		},
		{
			name: "test both",
			attrs: map[string]string{
				"field": "duration",
				"min":   "1s",
				"max":   "5h",
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
			filter, err := duration.NewDurationFilter(tt.attrs)
			require.NoError(t, err)
			if success := filter.Filter(context.Background(), &tt.silence); success != tt.expectedSuccess {
				t.Errorf("expected success: %t, but we didn't get it", tt.expectedSuccess)
			}
		})
	}
}
