package query

import (
	"context"
	"errors"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var ErrMissingType = errors.New("missing filter type")
var ErrStatsQueryFilterDoesntExist = errors.New("alert stats query filter doesn't exist")

// StatsResult is a single result from a StatsQuery.
type StatsResult struct {
	// Labels are the labels that apply to this result.
	Labels map[string]string

	// Frames are the data frames that apply to this result.
	Frames [][]float64
}

// AlertStatsQuery is a query that can be run against a DB to pull aggregated numbers out of it.
type AlertStatsQuery interface {
	// Filter returns the filter that this filters the data that gets passed to Process().
	Filter() AlertFilter

	// Process is called for each alert that matches the filter.
	Process(ctx context.Context, alert *model.Alert) error

	// Gather is called after all alerts have been processed. It returns the results of the query.
	Gather(ctx context.Context) []StatsResult
}

type alertStatsQueryConstructor func(args map[string]string) AlertStatsQuery

var alertStatsQueryRegistry = map[string]alertStatsQueryConstructor{}

// RegisterAlertStatsQuery registers a new AlertStatsQuery.
func RegisterAlertStatsQuery(name string, constructor alertStatsQueryConstructor) {
	alertStatsQueryRegistry[name] = constructor
}

// UnmarshalAlertStatsQuery unmarshals an AlertStatsQuery from a set of arguments.
func UnmarshalAlertStatsQuery(args map[string]string) (AlertStatsQuery, error) {
	name, ok := args["type"]
	if !ok {
		return nil, ErrMissingType
	}
	delete(args, "type")

	constructor, ok := alertStatsQueryRegistry[name]
	if !ok {
		return nil, ErrStatsQueryFilterDoesntExist
	}

	return constructor(args), nil
}

func init() {
	RegisterAlertStatsQuery("count", NewAlertCountQuery)
}

// AlertCountQuery counts the number of alerts that match the filter.
type AlertCountQuery struct {
	filter AlertFilter
	count  int
}

func NewAlertCountQuery(args map[string]string) AlertStatsQuery {
	if _, ok := args["filter_type"]; ok {
		filter, err := UnmarshalAlertFilter(args)
		if err != nil {
			return nil
		}

		return &AlertCountQuery{
			filter: filter,
		}
	}

	return &AlertCountQuery{}
}

func (q *AlertCountQuery) Filter() AlertFilter {
	return q.filter
}

func (q *AlertCountQuery) Process(ctx context.Context, alert *model.Alert) error {
	q.count += 1
	return nil
}

func (q *AlertCountQuery) Gather(ctx context.Context) []StatsResult {
	return []StatsResult{
		{
			Labels: map[string]string{},
			Frames: [][]float64{{float64(q.count)}},
		},
	}
}
