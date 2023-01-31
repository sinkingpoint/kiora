package kioradb

import (
	"context"

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

func (e *ExactLabelMatchQuery) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	if e.labelsHash == 0 {
		e.labelsHash = e.Labels.Hash()
	}

	return alert.Labels.Hash() == e.labelsHash
}
