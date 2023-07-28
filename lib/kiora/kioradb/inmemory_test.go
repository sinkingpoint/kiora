package kioradb_test

import (
	"context"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestInMemoryDB(t *testing.T) {
	db := kioradb.NewInMemoryDB()

	// Add an alert
	require.NoError(t, db.StoreAlerts(context.Background(), []model.Alert{
		{
			Labels: model.Labels{
				"foo": "bar",
			},
			Status: model.AlertStatusFiring,
		},
	}...))

	alerts := db.QueryAlerts(context.TODO(), query.NewAlertQuery(query.MatchAll()))

	require.Len(t, alerts, 1)
	alert := alerts[0]
	require.Equal(t, "bar", alert.Labels["foo"])
	require.Equal(t, model.AlertStatusFiring, alert.Status)

	// Resolve the above alert, add another
	require.NoError(t, db.StoreAlerts(context.Background(), []model.Alert{
		{
			Labels: model.Labels{
				"foo": "bar",
			},
			Status: model.AlertStatusResolved,
		},
		{
			Labels: model.Labels{
				"bar": "baz",
			},
			Status: model.AlertStatusFiring,
		},
	}...))

	alerts = db.QueryAlerts(context.TODO(), query.NewAlertQuery(query.MatchAll()))
	require.Len(t, alerts, 2)

	for _, alert := range alerts {
		if _, hasFoo := alert.Labels["foo"]; hasFoo {
			require.Equal(t, "bar", alert.Labels["foo"])
			require.Equal(t, model.AlertStatusResolved, alert.Status)
		} else if _, hasBar := alert.Labels["bar"]; hasBar {
			require.Equal(t, "baz", alert.Labels["bar"])
			require.Equal(t, model.AlertStatusFiring, alert.Status)
		} else {
			t.Errorf("unexpected alert: %v", alert)
		}
	}

	// Timeout the second alert
	require.NoError(t, db.StoreAlerts(context.Background(), []model.Alert{
		{
			Labels: model.Labels{
				"bar": "baz",
			},
			Status: model.AlertStatusResolved,
		},
	}...))

	alerts = db.QueryAlerts(context.TODO(), query.NewAlertQuery(query.MatchAll()))
	require.Len(t, alerts, 2)

	for _, alert := range alerts {
		if _, hasFoo := alert.Labels["foo"]; hasFoo {
			require.Equal(t, "bar", alert.Labels["foo"])
			require.Equal(t, model.AlertStatusResolved, alert.Status)
		} else if _, hasBar := alert.Labels["bar"]; hasBar {
			require.Equal(t, "baz", alert.Labels["bar"])
			require.Equal(t, model.AlertStatusResolved, alert.Status)
		} else {
			t.Errorf("unexpected alert: %v", alert)
		}
	}
}
