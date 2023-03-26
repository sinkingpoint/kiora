package timeout

import (
	"context"
	"time"

	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type TimeoutService struct {
	bus services.Bus
}

func NewTimeoutService(bus services.Bus) *TimeoutService {
	return &TimeoutService{
		bus: bus,
	}
}

func (t *TimeoutService) Name() string {
	return "timeout"
}

func (t *TimeoutService) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			t.timeoutAlerts(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (t *TimeoutService) timeoutAlerts(ctx context.Context) {
	query := query.NewAlertQuery(query.Status(model.AlertStatusFiring))
	changed := []model.Alert{}
	for _, a := range t.bus.DB().QueryAlerts(ctx, query) {
		if a.TimeOutDeadline.Before(time.Now()) {
			a.Status = model.AlertStatusTimedOut
			changed = append(changed, a)
		}
	}

	if err := t.bus.Broadcaster().BroadcastAlerts(ctx, changed...); err != nil {
		t.bus.Logger("timeout").Warn().Err(err).Msg("failed to broadcast alerts")
	}
}
