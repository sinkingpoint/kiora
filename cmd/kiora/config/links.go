package config

type Filter interface {
	Type() string
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
