package query

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

type Query struct {
	Order   Order
	OrderBy []string
	Offset  int
	Limit   int
}

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

func NewAlertQuery(filter AlertFilter, ops ...QueryOption) *AlertQuery {
	q := AlertQuery{
		Filter: filter,
	}

	for _, op := range ops {
		op.Apply(&q.Query)
	}

	return &q
}
