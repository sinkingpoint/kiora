package query_test

import (
	"sort"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
)

func TestSortAlertsByFields(t *testing.T) {
	a := model.Alert{
		StartTime: time.Unix(1, 0),
		Labels: model.Labels{
			"foo": "bar",
		},
		EndTime: time.Unix(2, 0),
		Status:  model.AlertStatusFiring,
	}

	b := model.Alert{
		StartTime: time.Unix(2, 0),
		Labels: model.Labels{
			"foo": "baz",
		},
		EndTime: time.Unix(2, 0),
		Status:  model.AlertStatusFiring,
	}

	c := model.Alert{
		StartTime: time.Unix(3, 0),
		Labels: model.Labels{
			"foo": "qux",
		},
		EndTime: time.Unix(2, 0),
		Status:  model.AlertStatusFiring,
	}

	tests := []struct {
		Name           string
		Alerts         []model.Alert
		Fields         []string
		Order          query.Order
		ExpectedAlerts []model.Alert
	}{
		{
			Name: "test_sort_by_start_time",
			Alerts: []model.Alert{
				a, c, b,
			},
			Fields: []string{"__starts_at__"},
			Order:  query.OrderAsc,
			ExpectedAlerts: []model.Alert{
				a, b, c,
			},
		},
		{
			Name: "test_sort_by_start_time_desc",
			Alerts: []model.Alert{
				a, c, b,
			},
			Fields: []string{"__starts_at__"},
			Order:  query.OrderDesc,
			ExpectedAlerts: []model.Alert{
				c, b, a,
			},
		},
		{
			Name: "test_sort_by_label",
			Alerts: []model.Alert{
				a, c, b,
			},
			Fields: []string{"foo"},
			Order:  query.OrderAsc,
			ExpectedAlerts: []model.Alert{
				a, b, c,
			},
		},
		{
			Name: "test_sort_by_multiple_values",
			Alerts: []model.Alert{
				a, c, b,
			},
			Fields: []string{"__ends_at__", "foo"},
			Order:  query.OrderDesc,
			ExpectedAlerts: []model.Alert{
				c, b, a,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sort.Sort(query.SortAlertsByFields(tt.Alerts, tt.Fields, tt.Order))
			for i, alert := range tt.Alerts {
				assert.Equal(t, tt.ExpectedAlerts[i], alert, "expected alert %d to be %v, got %v", i, tt.ExpectedAlerts[i], alert)
			}
		})
	}
}
