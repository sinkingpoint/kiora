package config

import (
	"context"
)

// Filter defines something that can filter models.
type Filter interface {
	// Type returns a string representation of the type of the filter.
	Type() string

	// Filter returns an error if the given model should be filtered.
	Filter(ctx context.Context, f Fielder) error
}

// Fielder is a thing that has fields that can be filtered.
type Fielder interface {
	// Field returns the value of a field.
	Field(name string) (any, error)

	// Fields returns a map of all the fields that can be filtered.
	Fields() map[string]any
}

// FilterConstructor is a function that can construct a filter from a set of attributes.
type FilterConstructor = func(globals *Globals, attrs map[string]string) (Filter, error)

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
