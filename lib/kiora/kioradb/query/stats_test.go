package query_test

import (
	"context"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestAlertCountQuery(t *testing.T) {
	tests := []struct {
		name   string
		ty     string
		args   map[string]string
		alerts []model.Alert
		want   float64
	}{
		{
			name: "test alert count query",
			ty:   "count",
			args: map[string]string{},
			alerts: []model.Alert{
				{
					Labels: model.Labels{
						"foo": "bar",
					},
				},
			},
			want: 1,
		},
		{
			name: "test alert count query with filter",
			ty:   "count",
			args: map[string]string{
				"filter_type": "status",
				"status":      "firing",
			},
			alerts: []model.Alert{
				{
					Labels: model.Labels{
						"foo": "bar",
					},
					Status: model.AlertStatusFiring,
				},
				{
					Labels: model.Labels{
						"foo": "baz",
					},
					Status: model.AlertStatusResolved,
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.UnmarshalAlertStatsQuery(tt.ty, tt.args)
			require.NoError(t, err, "unmarshal alert stats query failed")
			require.NotNil(t, q, "unmarshal alert stats query failed")

			filter := q.Filter()
			for _, alert := range tt.alerts {
				if filter == nil || filter.MatchesAlert(context.TODO(), &alert) {
					require.NoError(t, q.Process(context.TODO(), &alert), "process alert failed")
				}
			}

			got := q.Gather(context.TODO())
			require.Len(t, got, 1, "gather returned wrong number of results")
			require.Equal(t, tt.want, got[0].Frames[0][0], "gather returned wrong value")
		})
	}
}

func TestAlertStatusCountQuery(t *testing.T) {
	tests := []struct {
		name   string
		ty     string
		args   map[string]string
		alerts []model.Alert
		want   map[string]float64
	}{
		{
			name: "test alert status count query",
			ty:   "status_count",
			args: map[string]string{},
			alerts: []model.Alert{
				{
					Labels: model.Labels{
						"foo": "bar",
					},
					Status: model.AlertStatusFiring,
				},
				{
					Labels: model.Labels{
						"foo": "baz",
					},
					Status: model.AlertStatusFiring,
				},
				{
					Labels: model.Labels{
						"foo": "qux",
					},
					Status: model.AlertStatusResolved,
				},
			},
			want: map[string]float64{
				"firing":   2,
				"resolved": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.UnmarshalAlertStatsQuery(tt.ty, tt.args)
			require.NoError(t, err, "unmarshal alert stats query failed")
			require.NotNil(t, q, "unmarshal alert stats query failed")

			for _, alert := range tt.alerts {
				require.NoError(t, q.Process(context.TODO(), &alert), "process alert failed")
			}

			got := q.Gather(context.TODO())
			require.Len(t, got, len(tt.want), "gather returned wrong number of results")
			for _, result := range got {
				require.Equal(t, tt.want[result.Labels["status"]], result.Frames[0][0], "gather returned wrong value")
			}
		})
	}
}
