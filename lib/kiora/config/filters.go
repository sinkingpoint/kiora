package config

import (
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Filter defines something that can filter models.
type Filter interface {
	// Type returns a string representation of the type of the filter.
	Type() string
}

// AlertFilter is a type of filter that can filter alerts.
type AlertFilter interface {
	FilterAlert(a *model.Alert) bool
}

// AlertAcknowledgementFilter is a type of filter that can filter alert acknowledgements.
type AlertAcknowledgementFilter interface {
	FilterAlertAcknowledgement(alert *model.Alert, ack *model.AlertAcknowledgement) bool
}

// FilterConstructor is a function that can construct a filter from a set of attributes.
type FilterConstructor = func(attrs map[string]string) (Filter, error)

var filterRegistry = map[string]FilterConstructor{}

// LookupFilter looks up a filter constructor by name.
func LookupFilter(name string) (FilterConstructor, bool) {
	cons, ok := filterRegistry[name]
	return cons, ok
}

// RegisterFilter registers a filter constructor by name.
func RegisterFilter(name string, cons FilterConstructor) {
	filterRegistry[name] = cons
}
