package kioradb_test

import (
	"context"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryDB(t *testing.T) {
	db := kioradb.NewInMemoryDB()

	// Add an alert
	assert.NoError(t, db.StoreAlerts(context.Background(), []model.Alert{
		{
			Labels: model.Labels{
				"foo": "bar",
			},
			Status: model.AlertStatusFiring,
		},
	}...))

	alerts := db.QueryAlerts(context.TODO(), &kioradb.AllMatchQuery{})

	assert.Len(t, alerts, 1)
	alert := alerts[0]
	assert.Equal(t, "bar", alert.Labels["foo"])
	assert.Equal(t, model.AlertStatusFiring, alert.Status)

	// Resolve the above alert, add another
	assert.NoError(t, db.StoreAlerts(context.Background(), []model.Alert{
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

	alerts = db.QueryAlerts(context.TODO(), &kioradb.AllMatchQuery{})
	require.Len(t, alerts, 2)

	for _, alert := range alerts {
		if _, hasFoo := alert.Labels["foo"]; hasFoo {
			assert.Equal(t, "bar", alert.Labels["foo"])
			assert.Equal(t, model.AlertStatusResolved, alert.Status)
		} else if _, hasBar := alert.Labels["bar"]; hasBar {
			assert.Equal(t, "baz", alert.Labels["bar"])
			assert.Equal(t, model.AlertStatusFiring, alert.Status)
		} else {
			t.Errorf("unexpected alert: %q", alert)
		}
	}

	// Timeout the second alert
	assert.NoError(t, db.StoreAlerts(context.Background(), []model.Alert{
		{
			Labels: model.Labels{
				"bar": "baz",
			},
			Status: model.AlertStatusResolved,
		},
	}...))

	alerts = db.QueryAlerts(context.TODO(), &kioradb.AllMatchQuery{})
	require.Len(t, alerts, 2)

	for _, alert := range alerts {
		if _, hasFoo := alert.Labels["foo"]; hasFoo {
			assert.Equal(t, "bar", alert.Labels["foo"])
			assert.Equal(t, model.AlertStatusResolved, alert.Status)
		} else if _, hasBar := alert.Labels["bar"]; hasBar {
			assert.Equal(t, "baz", alert.Labels["bar"])
			assert.Equal(t, model.AlertStatusResolved, alert.Status)
		} else {
			t.Errorf("unexpected alert: %q", alert)
		}
	}
}
