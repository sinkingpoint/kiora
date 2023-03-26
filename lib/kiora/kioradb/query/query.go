package query

type Order string

type Query struct {
	Order   Order
	OrderBy []string
	Offset  int
	Limit   int
}

type QueryOp interface {
	Apply(*Query)
}

type QueryOpFunc func(*Query)

func (f QueryOpFunc) Apply(q *Query) {
	f(q)
}

func OrderBy(fields []string, order Order) QueryOp {
	return QueryOpFunc(func(q *Query) {
		q.OrderBy = fields
		q.Order = order
	})
}

func Limit(limit int) QueryOp {
	return QueryOpFunc(func(q *Query) {
		q.Limit = limit
	})
}

func Offset(offset int) QueryOp {
	return QueryOpFunc(func(q *Query) {
		q.Offset = offset
	})
}

type AlertQuery struct {
	Query
	Filter AlertFilter
}

func NewAlertQuery(filter AlertFilter, ops ...QueryOp) *AlertQuery {
	q := AlertQuery{
		Filter: filter,
	}

	for _, op := range ops {
		op.Apply(&q.Query)
	}

	return &q
}
