package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
)

// NewFilter implements config.FilterFactory.
func NewFilter(globals *config.Globals, attrs map[string]string) (config.Filter, error) {
	delete(attrs, "type")
	rateLimitFilter := RateLimitFilter{
		globals: globals,
	}

	if err := unmarshal.UnmarshalConfig(attrs, &rateLimitFilter, unmarshal.UnmarshalOpts{DisallowUnknownFields: true}); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal rate limit filter")
	}

	// If burst is not set, default to rate.
	if rateLimitFilter.Burst == 0 {
		rateLimitFilter.Burst = rateLimitFilter.Rate
	}

	return &rateLimitFilter, nil
}

// RateLimitFilter implements a tenant-aware rate limit filter.
type RateLimitFilter struct {
	globals *config.Globals

	// Interval is the time interval over which the rate limit applies.
	Interval time.Duration `config:"interval" required:"true"`

	// Rate is the number of alerts allowed per interval.
	Rate int `config:"rate" required:"true"`

	// Burst is the number of alerts allowed to exceed the rate limit.
	Burst int `config:"burst"`

	buckets sync.Map
}

func (r *RateLimitFilter) newBucket() *ratelimitBucket {
	return &ratelimitBucket{
		lock:       sync.Mutex{},
		interval:   r.Interval,
		rate:       r.Rate,
		burst:      r.Burst,
		tokenCount: r.Rate,
		lastUpdate: time.Now(),
	}
}

// Filter implements config.Filter.
func (r *RateLimitFilter) Filter(ctx context.Context, f config.Fielder) error {
	tenant, err := r.globals.Tenanter.GetTenant(ctx, f)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	bucket, _ := r.buckets.LoadOrStore(tenant, r.newBucket())

	if !bucket.(*ratelimitBucket).consumeToken() {
		return fmt.Errorf("rate limit of %d per %s exceeded for tenant %s", r.Rate, r.Interval, tenant)
	}

	return nil
}

// Type implements config.Filter.
func (*RateLimitFilter) Type() string {
	return "rate_limit"
}

// ratelimitBucket implements a token bucket rate limiter.
type ratelimitBucket struct {
	lock       sync.Mutex
	tokenCount int
	lastUpdate time.Time

	interval time.Duration
	rate     int
	burst    int
}

// updateCount updates the token count based on the time since the last update.
func (b *ratelimitBucket) updateCount() {
	timeSinceLastUpdate := time.Since(b.lastUpdate)
	newTokens := float64(timeSinceLastUpdate) / float64(b.interval) * float64(b.rate)
	if newTokens > 0 {
		b.tokenCount += int(newTokens)
		b.lastUpdate = stubs.Time.Now()
		if b.tokenCount > b.burst {
			b.tokenCount = b.burst
		}
	}
}

func (b *ratelimitBucket) consumeToken() bool {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.updateCount()

	if b.tokenCount > 0 {
		b.tokenCount--
		return true
	}

	return false
}
