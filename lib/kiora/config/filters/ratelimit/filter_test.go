package ratelimit_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/ratelimit"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	globals := config.NewGlobals(config.WithTenanter(config.NewStaticTenanter("")))
	filter, err := ratelimit.NewFilter(globals, map[string]string{
		"interval": "1s",
		"rate":     "1",
		"burst":    "2",
	})

	require.NoError(t, err)

	alert := model.Alert{
		Status: model.AlertStatusFiring,
		Labels: model.Labels{},
	}

	require.NoError(t, alert.Materialise())

	// Assert that we can send one alert, but the second exceeds rate limits and fails.
	require.NoError(t, filter.Filter(context.Background(), &alert))
	require.Error(t, filter.Filter(context.Background(), &alert))

	// Sleep for a bit to allow the interval to pass.
	time.Sleep(2 * time.Second)

	// Assert that we can send two alerts after the interval has passed, because of the burst capacity.
	require.NoError(t, filter.Filter(context.Background(), &alert))
	require.NoError(t, filter.Filter(context.Background(), &alert))
}

// TestRatelimitConcurrency tests that the ratelimit filter is safe to use concurrently,
// and that it correctly enforces the rate limit.
func TestRatelimitConcurrency(t *testing.T) {
	numSuccess := atomic.Int32{}

	globals := config.NewGlobals(config.WithTenanter(config.NewStaticTenanter("")))
	filter, err := ratelimit.NewFilter(globals, map[string]string{
		"interval": "30s",
		"rate":     "200",
	})

	require.NoError(t, err)

	alert := model.Alert{
		Status: model.AlertStatusFiring,
		Labels: model.Labels{},
	}

	require.NoError(t, alert.Materialise())

	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			err := filter.Filter(context.Background(), &alert)
			if err == nil {
				numSuccess.Add(1)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	require.Equal(t, 200, int(numSuccess.Load()))
}
