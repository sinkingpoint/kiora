package query

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

// Query represents a generic query based on SQL semantics.
type Query struct {
	// OrderBy is a direction to order the results by.
	Order Order

	// OrderBy is a list of fields to order the results by.
	OrderBy []string

	// Offset is the number of results to skip before returning results.
	Offset int

	// Limit is the maximum number of results to return.
	Limit int
}

// QueryOption represents an option that can be applied to a query.
type QueryOption interface {
	Apply(*Query)
}

type QueryOpFunc func(*Query)

func (f QueryOpFunc) Apply(q *Query) {
	f(q)
}

func OrderBy(fields []string, order Order) QueryOption {
	return QueryOpFunc(func(q *Query) {
		q.OrderBy = fields
		q.Order = order
	})
}

func Limit(limit int) QueryOption {
	return QueryOpFunc(func(q *Query) {
		q.Limit = limit
	})
}

func Offset(offset int) QueryOption {
	return QueryOpFunc(func(q *Query) {
		q.Offset = offset
	})
}

type AlertQuery struct {
	Query
	Filter AlertFilter
}

func NewAlertQuery(filter AlertFilter, ops ...QueryOption) AlertQuery {
	q := AlertQuery{
		Filter: filter,
	}

	for _, op := range ops {
		op.Apply(&q.Query)
	}

	return q
}

type SilenceQuery struct {
	Query
	Filter SilenceFilter
}

func NewSilenceQuery(filter SilenceFilter, ops ...QueryOption) SilenceQuery {
	q := SilenceQuery{
		Filter: filter,
	}

	for _, op := range ops {
		op.Apply(&q.Query)
	}

	return q
}
