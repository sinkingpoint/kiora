package config

import "github.com/sinkingpoint/kiora/lib/kiora/model"

type Filter interface {
	Type() string
	FilterAlert(a *model.Alert) bool
}

type Link struct {
	incomingFilter Filter
	to             string
}

type FilterConstructor = func(n edge) (Filter, error)

var filterRegistry = map[string]FilterConstructor{}

func LookupFilter(name string) FilterConstructor {
	return filterRegistry[name]
}
