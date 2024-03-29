package query

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// AlertFilterConstructor is a function that can construct an AlertFilter from a set of arguments.
type AlertFilterConstructor func(args map[string]string) AlertFilter

var alertFilterRegistry = map[string]AlertFilterConstructor{}

func RegisterAlertFilter(name string, constructor AlertFilterConstructor) {
	alertFilterRegistry[name] = constructor
}

// UnmarshalAlertFilter unmarshals an AlertFilter from a set of arguments.
func UnmarshalAlertFilter(args map[string]string) (AlertFilter, error) {
	name, ok := args["filter_type"]
	if !ok {
		return nil, errors.New("missing filter type")
	}
	delete(args, "filter_type")

	constructor, ok := alertFilterRegistry[name]
	if !ok {
		return nil, fmt.Errorf("unknown filter type %q", name)
	}

	return constructor(args), nil
}

func init() {
	RegisterAlertFilter("exact", func(args map[string]string) AlertFilter {
		return ExactLabelMatch(model.Labels(args))
	})

	RegisterAlertFilter("partial", func(args map[string]string) AlertFilter {
		return PartialLabelMatch(model.Labels(args))
	})

	RegisterAlertFilter("status", func(args map[string]string) AlertFilter {
		status, ok := args["status"]
		if !ok {
			return nil
		}
		return Status(model.AlertStatus(status))
	})

	RegisterAlertFilter("all", func(args map[string]string) AlertFilter {
		return MatchAll()
	})
}

// AlertFilter is a query that can be run against a DB to pull things out of it.
type AlertFilter interface {
	Type() string
	MatchesAlert(ctx context.Context, alert *model.Alert) bool
}

type AlertFilterFunc func(ctx context.Context, alert *model.Alert) bool

func (a AlertFilterFunc) Type() string {
	return "func"
}

func (a AlertFilterFunc) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return a(ctx, alert)
}

// SilenceFilter is a query that can be run against a DB to pull things out of it.
type SilenceFilter interface {
	MatchesSilence(ctx context.Context, silence *model.Silence) bool
}

type SilenceFilterFunc func(ctx context.Context, alert *model.Silence) bool

func (a SilenceFilterFunc) MatchesSilence(ctx context.Context, alert *model.Silence) bool {
	return a(ctx, alert)
}

// PartialLabelMatchFilter is an AlertFilter that matches alerts that contain the given labels (but may have extras on top of these).
type PartialLabelMatchFilter struct {
	Labels model.Labels
}

func PartialLabelMatch(labels model.Labels) *PartialLabelMatchFilter {
	return &PartialLabelMatchFilter{
		Labels: labels,
	}
}

func (p *PartialLabelMatchFilter) Type() string {
	return "partial"
}

func (p *PartialLabelMatchFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	for k, v := range p.Labels {
		if alertV, ok := alert.Labels[k]; !ok || alertV != v {
			return false
		}
	}

	return true
}

func (p *PartialLabelMatchFilter) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	return silence.Matches(p.Labels)
}

// ExactLabelMatchFilter is an AlertFilter that matches alerts that contain exactly the given labelset.
type ExactLabelMatchFilter struct {
	Labels     model.Labels
	labelsHash model.LabelsHash
}

func ExactLabelMatch(labels model.Labels) AlertFilter {
	return &ExactLabelMatchFilter{
		Labels: labels,
	}
}

func (e *ExactLabelMatchFilter) Type() string {
	return "exact"
}

func (e *ExactLabelMatchFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	if e.labelsHash == 0 {
		e.labelsHash = e.Labels.Hash()
	}

	return alert.Labels.Hash() == e.labelsHash
}

// AllMatchFilter is a query that returns every alert.
type AllMatchFilter struct{}

func MatchAll() *AllMatchFilter {
	return &AllMatchFilter{}
}

func (a *AllMatchFilter) Type() string {
	return "all"
}

func (a *AllMatchFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return true
}

func (a *AllMatchFilter) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	return true
}

// StatusFilter returns all the alerts that match the given status.
type StatusFilter struct {
	Status model.AlertStatus
}

func Status(s model.AlertStatus) *StatusFilter {
	return &StatusFilter{
		Status: s,
	}
}

func (s *StatusFilter) Type() string {
	return "status"
}

func (s *StatusFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return alert.Status == s.Status
}

type LastNotifyTimeRangeFilter struct {
	MinTime time.Time
	MaxTime time.Time
}

func (l *LastNotifyTimeRangeFilter) Type() string {
	return "last_notify_time"
}

func (l *LastNotifyTimeRangeFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	if !l.MinTime.IsZero() && alert.LastNotifyTime.Before(l.MinTime) {
		return false
	}

	if !l.MaxTime.IsZero() && alert.LastNotifyTime.After(l.MaxTime) {
		return false
	}

	return true
}

func LastNotifyTimeMin(minTime time.Time) *LastNotifyTimeRangeFilter {
	return &LastNotifyTimeRangeFilter{
		MinTime: minTime,
	}
}

func LastNotifyTimeMax(maxTime time.Time) *LastNotifyTimeRangeFilter {
	return &LastNotifyTimeRangeFilter{
		MaxTime: maxTime,
	}
}

func LastNotifyTimeWithin(minTime, maxTime time.Time) *LastNotifyTimeRangeFilter {
	return &LastNotifyTimeRangeFilter{
		MinTime: minTime,
		MaxTime: maxTime,
	}
}

type AllFilter struct {
	alertQueries   []AlertFilter
	silenceQueries []SilenceFilter
}

func (a *AllFilter) Type() string {
	alertTypes := []string{}
	for _, q := range a.alertQueries {
		alertTypes = append(alertTypes, q.Type())
	}

	return fmt.Sprintf("all(%s)", strings.Join(alertTypes, ","))
}

func (a *AllFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	for _, q := range a.alertQueries {
		if !q.MatchesAlert(ctx, alert) {
			return false
		}
	}

	return true
}

func (a *AllFilter) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	for _, q := range a.silenceQueries {
		if !q.MatchesSilence(ctx, silence) {
			return false
		}
	}

	return true
}

func AllAlerts(queries ...AlertFilter) *AllFilter {
	return &AllFilter{
		alertQueries: queries,
	}
}

func AllSilences(queries ...SilenceFilter) *AllFilter {
	return &AllFilter{
		silenceQueries: queries,
	}
}

// IDFilter is a query that matches a specific alert or silence by ID.
type IDFilter struct {
	ID string
}

func (i *IDFilter) Type() string {
	return "id"
}

func (i *IDFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return alert.ID == i.ID
}

func (i *IDFilter) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	return silence.ID == i.ID
}

func ID(id string) *IDFilter {
	return &IDFilter{
		ID: id,
	}
}

// SilenceIsActive returns a SilenceFilter that matches only active silences.
func SilenceIsActive() SilenceFilter {
	return SilenceFilterFunc(func(ctx context.Context, silence *model.Silence) bool {
		return silence.IsActive()
	})
}

// MatcherFilter is a Filter that matches alerts or silences that contain the given matcher.
type MatcherFilter struct {
	matcher model.Matcher
}

func Matcher(m model.Matcher) *MatcherFilter {
	return &MatcherFilter{
		matcher: m,
	}
}

func (i *MatcherFilter) Type() string {
	return "matcher"
}

// MatchesAlert returns true if the given matcher matcher the given alert.
func (m *MatcherFilter) MatchesAlert(ctx context.Context, alert *model.Alert) bool {
	return m.matcher.Matches(alert.Labels)
}

// MatchesSilence returns true if the given matcher is in the given silence.
func (m *MatcherFilter) MatchesSilence(ctx context.Context, silence *model.Silence) bool {
	for i := range silence.Matchers {
		matcher := silence.Matchers[i]
		if matcher.Label == m.matcher.Label && matcher.Value == m.matcher.Value && matcher.IsNegative == m.matcher.IsNegative && matcher.IsRegex == m.matcher.IsRegex {
			return true
		}
	}

	return false
}
