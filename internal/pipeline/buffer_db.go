package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel"
)

// bufferDB is a database that buffers alerts and silences in memory and periodically flushes them to the underlying database.
type bufferDB struct {
	// aLock protects the alertsBuffer.
	aLock        sync.Mutex
	alertsBuffer []model.Alert

	// sLock protects the silencesBuffer.
	sLock          sync.Mutex
	silencesBuffer []model.Silence

	// lengthLimit is the maximum number of alerts or silences that can be buffered before they are flushed.
	lengthLimit int

	// timeLimit is the maximum amount of time that alerts or silences can be buffered before they are flushed.
	// It's recommended to keep this low to ensure that alerts and silences are flushed quickly.
	timeLimit time.Duration

	// bus is the bus that the bufferDB is attached to.
	bus services.Bus
}

// NewBufferDB creates a new bufferDB.
func NewBufferDB(bus services.Bus, alertCapacity int, silenceCapacity int, lengthLimit int, timeLimit time.Duration) *bufferDB {
	return &bufferDB{
		lengthLimit:    lengthLimit,
		timeLimit:      timeLimit,
		bus:            bus,
		alertsBuffer:   make([]model.Alert, 0, alertCapacity),
		silencesBuffer: make([]model.Silence, 0, silenceCapacity),
	}
}

// Run is responsible for periodically flushing the buffer.
func (n *bufferDB) Run(ctx context.Context) error {
	ticker := time.NewTicker(n.timeLimit)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return n.Flush(ctx)
		case <-ticker.C:
			n.Flush(ctx)
		}
	}
}

// StoreAlerts stores the given alerts in the buffer, flushing the buffer if it exceeds the length limit.
func (b *bufferDB) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.StoreAlerts")
	defer span.End()

	b.aLock.Lock()
	defer b.aLock.Unlock()

	for len(b.alertsBuffer)+len(alerts) > b.lengthLimit {
		amt := b.lengthLimit - len(b.alertsBuffer)
		b.alertsBuffer = append(b.alertsBuffer, alerts[:amt]...)
		b.flushAlerts(ctx)
		alerts = alerts[amt:]
	}

	b.alertsBuffer = append(b.alertsBuffer, alerts...)
	return nil
}

// StoreSilences stores the given silences in the buffer, flushing the buffer if it exceeds the length limit.
func (b *bufferDB) StoreSilences(ctx context.Context, silences ...model.Silence) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.StoreSilences")
	defer span.End()

	b.sLock.Lock()
	defer b.sLock.Unlock()

	for len(b.silencesBuffer)+len(silences) > b.lengthLimit {
		amt := b.lengthLimit - len(b.silencesBuffer)
		b.silencesBuffer = append(b.silencesBuffer, silences[:amt]...)
		b.flushSilences(ctx)
		silences = silences[amt:]
	}

	b.silencesBuffer = append(b.silencesBuffer, silences...)
	return nil
}

func (b *bufferDB) flushAlerts(ctx context.Context) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.flushAlerts")
	defer span.End()

	b.aLock.Lock()
	defer b.aLock.Unlock()
	if err := b.bus.DB().StoreAlerts(ctx, b.alertsBuffer...); err != nil {
		return err
	}
	b.alertsBuffer = b.alertsBuffer[:0]
	return nil
}

func (b *bufferDB) flushSilences(ctx context.Context) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.flushSilences")
	defer span.End()

	b.sLock.Lock()
	defer b.sLock.Unlock()
	if err := b.bus.DB().StoreSilences(ctx, b.silencesBuffer...); err != nil {
		return err
	}
	b.silencesBuffer = b.silencesBuffer[:0]
	return nil
}

// Flush flushes the buffer to the underlying database.
func (b *bufferDB) Flush(ctx context.Context) error {
	var flushErr error
	if err := b.flushAlerts(ctx); err != nil {
		err = multierror.Append(flushErr, err)
	}

	if err := b.flushSilences(ctx); err != nil {
		err = multierror.Append(flushErr, err)
	}

	return flushErr
}
