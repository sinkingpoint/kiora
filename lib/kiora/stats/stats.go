package stats

import (
	"context"

	"github.com/sinkingpoint/kiora/internal/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Stats provides an interface to query statistics about models in the database.
type Stats struct {
	bus services.Bus
}

func NewStats(bus services.Bus) *Stats {
	return &Stats{
		bus: bus,
	}
}

// ExecuteAlertFilter executes an AlertFilter and returns the results.
func (s *Stats) ExecuteAlertFilter(ctx context.Context, aq AlertFilter) ([]StatsResult, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Stats.ExecuteAlertFilter")
	defer span.End()

	span.SetAttributes(attribute.String("query", aq.Name()))

	q := aq.Query(ctx)
	for _, alert := range s.bus.DB().QueryAlerts(ctx, q) {
		if err := aq.Ingest(ctx, &alert); err != nil {
			return nil, err
		}
	}

	return aq.Gather(ctx)
}
