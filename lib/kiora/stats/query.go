package stats

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func init() {
	RegisterAlertQuery("AlertStatusQuery", AlertStatusQuery)
}

var _ = AlertQuery(&alertStatusQuery{})

// StatsResult is a single result from a StatsQuery.
type StatsResult struct {
	Key   any
	Value int
}

// Query is an interface for querying statistics about models in the database.
type Query interface {
	Name() string
	// Gather is called after all the models have been sent through, to return the final results.
	Gather(ctx context.Context) ([]StatsResult, error)
}

// AlertQuery is an interface for querying statistics about alerts in the database.
type AlertQuery interface {
	Query

	// Query returns a query to use to find the alerts to send through the Ingest method.
	Query(ctx context.Context) query.AlertQuery

	// Ingest is called for each alert that matches the query. It should update the internal state of the query.
	Ingest(ctx context.Context, f *model.Alert) error
}

// alertStatusQuery is an AlertQuery that counts the number of alerts with each status.
type alertStatusQuery struct {
	counts map[model.AlertStatus]int
}

func AlertStatusQuery(attrs map[string]string) AlertQuery {
	return &alertStatusQuery{
		counts: make(map[model.AlertStatus]int),
	}
}

func (q *alertStatusQuery) Name() string {
	return "AlertStatusQuery"
}

func (q *alertStatusQuery) Gather(ctx context.Context) ([]StatsResult, error) {
	var results []StatsResult
	for k, v := range q.counts {
		results = append(results, StatsResult{Key: k, Value: v})
	}
	return results, nil
}

func (q *alertStatusQuery) Query(ctx context.Context) query.AlertQuery {
	return query.MatchAll()
}

func (q *alertStatusQuery) Ingest(ctx context.Context, f *model.Alert) error {
	q.counts[f.Status]++
	return nil
}
