package config

import "github.com/sinkingpoint/kiora/lib/kiora/model"

type Filter interface {
	Type() string
	FilterAlert(a *model.Alert) bool
}

type FilterConstructor = func(attrs map[string]string) (Filter, error)

var filterRegistry = map[string]FilterConstructor{}

func LookupFilter(name string) (FilterConstructor, bool) {
	cons, ok := filterRegistry[name]
	return cons, ok
}

func RegisterFilter(name string, cons FilterConstructor) {
	filterRegistry[name] = cons
}
