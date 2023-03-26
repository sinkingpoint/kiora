package query

import (
	"sort"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// alertsByFields implements sort.Interface for sorting a number of alerts
// by a number of fields. If a field is not present on an alert, it is considered
// to be less than the other alert.
type alertsByFields struct {
	Alerts []model.Alert
	Fields []string
	Order  Order
}

func SortAlertsByFields(alerts []model.Alert, fields []string, order Order) sort.Interface {
	return alertsByFields{
		Alerts: alerts,
		Fields: fields,
		Order:  order,
	}
}

func (a alertsByFields) Len() int {
	return len(a.Alerts)
}

func (a alertsByFields) Swap(i, j int) {
	a.Alerts[i], a.Alerts[j] = a.Alerts[j], a.Alerts[i]
}

func (a alertsByFields) lessVal() bool {
	return a.Order == OrderAsc
}

func (a alertsByFields) Less(i, j int) bool {
	for _, field := range a.Fields {
		iVal, iErr := a.Alerts[i].Field(field)
		jVal, jErr := a.Alerts[j].Field(field)

		if iErr != nil && jErr != nil {
			continue
		}

		if iErr != nil {
			return !a.lessVal()
		}

		if jErr != nil {
			return a.lessVal()
		}

		if iVal == jVal {
			continue
		}

		switch iVal.(type) {
		case string:
			if a.Order == OrderDesc {
				return (iVal.(string) > jVal.(string))
			}

			return iVal.(string) < jVal.(string)
		case int:
			if a.Order == OrderDesc {
				return (iVal.(int) > jVal.(int))
			}

			return iVal.(int) < jVal.(int)
		case float64:
			if a.Order == OrderDesc {
				return (iVal.(float64) > jVal.(float64))
			}

			return iVal.(float64) < jVal.(float64)
		case time.Time:
			if a.Order == OrderDesc {
				return (iVal.(time.Time).After(jVal.(time.Time)))
			}

			return iVal.(time.Time).Before(jVal.(time.Time))
		default:
			log.Warn().Str("field", field).Interface("value", iVal).Msg("unknown field type")
			continue
		}
	}

	return true
}
