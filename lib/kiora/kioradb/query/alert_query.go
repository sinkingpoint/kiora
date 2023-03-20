package query

import (
	"context"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// AlertQuery is a query that can be run against a DB to pull things out of it.
type AlertQuery interface {
	MatchesAlert(ctx context.Context, alert *model.Alert) bool
}

type AlertQueryFunc func(ctx context.Context, alert *model.Alert) bool

func (a AlertQueryFunc) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return a(ctx, alert)
}

// SilenceQuery is a query that can be run against a DB to pull things out of it.
type SilenceQuery interface {
	MatchesSilence(ctx context.Context, silence *model.Silence) bool
}

type SilenceQueryFunc func(ctx context.Context, alert *model.Silence) bool

func (a SilenceQueryFunc) MatchesSilence(ctx context.Context, alert *model.Silence) bool {
	return a(ctx, alert)
}

// PartialLabelMatchQuery is an AlertQuery that matches alerts that contain the given labels (but may have extras on top of these).
type PartialLabelMatchQuery struct {
	Labels model.Labels
}

func PartialLabelMatch(labels model.Labels) *PartialLabelMatchQuery {
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

func (p *PartialLabelMatchQuery) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	return silence.Matches(p.Labels)
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

func MatchAll() *AllMatchQuery {
	return &AllMatchQuery{}
}

func (a *AllMatchQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return true
}

func (a *AllMatchQuery) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	return true
}

// StatusQuery returns all the alerts that match the given status.
type StatusQuery struct {
	Status model.AlertStatus
}

func Status(s model.AlertStatus) *StatusQuery {
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
	alertQueries   []AlertQuery
	silenceQueries []SilenceQuery
}

func (a *AllQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	for _, q := range a.alertQueries {
		if !q.MatchesAlert(ctx, alert) {
			return false
		}
	}

	return true
}

func (a *AllQuery) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	for _, q := range a.silenceQueries {
		if !q.MatchesSilence(ctx, silence) {
			return false
		}
	}

	return true
}

func AllAlerts(queries ...AlertQuery) *AllQuery {
	return &AllQuery{
		alertQueries: queries,
	}
}

func AllSilences(queries ...SilenceQuery) *AllQuery {
	return &AllQuery{
		silenceQueries: queries,
	}
}

// IDQuery is a query that matches a specific alert by ID.
type IDQuery struct {
	ID string
}

func (i *IDQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return alert.ID == i.ID
}

func (i *IDQuery) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	return silence.ID == i.ID
}

func ID(id string) *IDQuery {
	return &IDQuery{
		ID: id,
	}
}

func SilenceIsActive() SilenceQuery {
	return SilenceQueryFunc(func(ctx context.Context, silence *model.Silence) bool {
		return silence.IsActive()
	})
}
