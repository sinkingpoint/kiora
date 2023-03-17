package query

import (
	"context"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Query is a query that can be run against a DB to pull things out of it.
type AlertQuery interface {
	MatchesAlert(ctx context.Context, alert *model.Alert) bool
}

// AlertQueryFunc provides a wrapper around an anonymous function to make it into an AlertQuery.
type AlertQueryFunc func(ctx context.Context, alert *model.Alert) bool

func (a AlertQueryFunc) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	if a != nil {
		return a(ctx, alert)
	}

	return false
}

// PartialLabelMatchQuery is an AlertQuery that matches alerts that contain the given labels (but may have extras on top of these).
type PartialLabelMatchQuery struct {
	Labels model.Labels
}

func PartialLabelMatch(labels model.Labels) AlertQuery {
	return &PartialLabelMatchQuery{
		Labels: labels,
	}
}

func (p *PartialLabelMatchQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	for k, v := range p.Labels {
		if alertV, ok := alert.Labels[k]; !ok || alertV != v {
			return false
		}
	}

	return true
}

// ExactLabelMatchQuery is an AlertQuery that matches alerts that contain exactly the given labelset.
type ExactLabelMatchQuery struct {
	Labels     model.Labels
	labelsHash model.LabelsHash
}

func ExactLabelMatch(labels model.Labels) AlertQuery {
	return &ExactLabelMatchQuery{
		Labels: labels,
	}
}

func (e *ExactLabelMatchQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	if e.labelsHash == 0 {
		e.labelsHash = e.Labels.Hash()
	}

	return alert.Labels.Hash() == e.labelsHash
}

// AllMatchQuery is a query that returns every alert.
type AllMatchQuery struct {
}

func MatchAll() AlertQuery {
	return &AllMatchQuery{}
}

func (a *AllMatchQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return true
}

// StatusQuery returns all the alerts that match the given status.
type StatusQuery struct {
	Status model.AlertStatus
}

func Status(s model.AlertStatus) AlertQuery {
	return &StatusQuery{
		Status: s,
	}
}

func (s *StatusQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return alert.Status == s.Status
}

type LastNotifyTimeRangeQuery struct {
	MinTime time.Time
	MaxTime time.Time
}

func (l *LastNotifyTimeRangeQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	if !l.MinTime.IsZero() && alert.LastNotifyTime.Before(l.MinTime) {
		return false
	}

	if !l.MaxTime.IsZero() && alert.LastNotifyTime.After(l.MaxTime) {
		return false
	}

	return true
}

func LastNotifyTimeMin(minTime time.Time) *LastNotifyTimeRangeQuery {
	return &LastNotifyTimeRangeQuery{
		MinTime: minTime,
	}
}

func LastNotifyTimeMax(maxTime time.Time) *LastNotifyTimeRangeQuery {
	return &LastNotifyTimeRangeQuery{
		MaxTime: maxTime,
	}
}

func LastNotifyTimeWithin(minTime, maxTime time.Time) *LastNotifyTimeRangeQuery {
	return &LastNotifyTimeRangeQuery{
		MinTime: minTime,
		MaxTime: maxTime,
	}
}

type AllQuery struct {
	queries []AlertQuery
}

func (a *AllQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	for _, q := range a.queries {
		if !q.MatchesAlert(ctx, alert) {
			return false
		}
	}

	return true
}

func All(queries ...AlertQuery) *AllQuery {
	return &AllQuery{
		queries: queries,
	}
}
