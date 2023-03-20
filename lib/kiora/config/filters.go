package config

import "github.com/sinkingpoint/kiora/lib/kiora/model"

type Filter interface {
	Type() string
}

type AlertFilter interface {
	FilterAlert(a *model.Alert) bool
}

type AlertAcknowledgementFilter interface {
	FilterAlertAcknowledgement(alert *model.Alert, ack *model.AlertAcknowledgement) bool
}

type FilterConstructor = func(attrs map[string]string) (Filter, error)

var filterRegistry = map[string]FilterConstructor{
	"": func(attrs map[string]string) (Filter, error) { return &NopFilter{}, nil },
}

func LookupFilter(name string) (FilterConstructor, bool) {
	cons, ok := filterRegistry[name]
	return cons, ok
}

func RegisterFilter(name string, cons FilterConstructor) {
	filterRegistry[name] = cons
}

type NopFilter struct{}

func (n *NopFilter) Type() string {
	return "nop"
}
