package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
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

	db kioradb.DB
}

// NewBufferDB creates a new bufferDB.
func NewBufferDB(db kioradb.DB, alertCapacity, silenceCapacity, lengthLimit int, timeLimit time.Duration) *bufferDB {
	return &bufferDB{
		lengthLimit:    lengthLimit,
		timeLimit:      timeLimit,
		db:             db,
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

	var flushErr error

	for len(b.alertsBuffer)+len(alerts) > b.lengthLimit {
		amt := b.lengthLimit - len(b.alertsBuffer)
		b.alertsBuffer = append(b.alertsBuffer, alerts[:amt]...)
		if err := b.flushAlerts(ctx); err != nil {
			flushErr = multierror.Append(flushErr, err)
		}

		alerts = alerts[amt:]
	}

	b.alertsBuffer = append(b.alertsBuffer, alerts...)
	return flushErr
}

// StoreSilences stores the given silences in the buffer, flushing the buffer if it exceeds the length limit.
func (b *bufferDB) StoreSilences(ctx context.Context, silences ...model.Silence) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.StoreSilences")
	defer span.End()

	b.sLock.Lock()
	defer b.sLock.Unlock()

	var flushErr error

	for len(b.silencesBuffer)+len(silences) > b.lengthLimit {
		amt := b.lengthLimit - len(b.silencesBuffer)
		b.silencesBuffer = append(b.silencesBuffer, silences[:amt]...)
		if err := b.flushSilences(ctx); err != nil {
			flushErr = multierror.Append(flushErr, err)
		}

		silences = silences[amt:]
	}

	b.silencesBuffer = append(b.silencesBuffer, silences...)
	return nil
}

// flushAlerts flushes the alerts buffer to the underlying database.
// NOTE: This function is not thread-safe and should only be called if the aLock is held.
func (b *bufferDB) flushAlerts(ctx context.Context) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.flushAlerts")
	defer span.End()

	if err := b.db.StoreAlerts(ctx, b.alertsBuffer...); err != nil {
		return errors.Wrap(err, "failed to store alerts")
	}
	b.alertsBuffer = b.alertsBuffer[:0]
	return nil
}

// flushSilences flushes the silences buffer to the underlying database.
// NOTE: This function is not thread-safe and should only be called if the sLock is held.
func (b *bufferDB) flushSilences(ctx context.Context) error {
	ctx, span := otel.Tracer("").Start(ctx, "bufferDB.flushSilences")
	defer span.End()

	if err := b.db.StoreSilences(ctx, b.silencesBuffer...); err != nil {
		return err
	}
	b.silencesBuffer = b.silencesBuffer[:0]
	return nil
}

// Flush flushes the buffer to the underlying database.
func (b *bufferDB) Flush(ctx context.Context) error {
	var flushErr error

	b.aLock.Lock()
	if err := b.flushAlerts(ctx); err != nil {
		flushErr = multierror.Append(flushErr, err)
	}
	b.aLock.Unlock()

	b.sLock.Lock()
	if err := b.flushSilences(ctx); err != nil {
		flushErr = multierror.Append(flushErr, err)
	}

	b.sLock.Unlock()

	return flushErr
}
