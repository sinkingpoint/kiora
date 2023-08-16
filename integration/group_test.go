package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestGrouping(t *testing.T) {
	initT(t)

	k := NewKioraInstance(t).WithConfig(`digraph config {
		wait_1s [type="group_wait" duration="1s"];
		group_by_alertname [type="group_labels" labels="alertname"];
		console [type="stdout"];

		alerts -> group_by_alertname -> wait_1s -> console;
	}`).Start()

	defer func() {
		require.NoError(t, k.Stop())
	}()

	alert1 := model.Alert{
		Labels: model.Labels{
			"alertname": "test",
			"foo":       "bar",
		},
		Annotations: make(map[string]string),
		Status:      model.AlertStatusFiring,
	}

	require.NoError(t, alert1.Materialise())

	alert2 := model.Alert{
		Labels: model.Labels{
			"alertname": "test",
			"foo":       "baz",
		},
		Annotations: make(map[string]string),
		Status:      model.AlertStatusFiring,
	}

	require.NoError(t, alert2.Materialise())

	k.SendAlert(context.TODO(), alert1)
	k.SendAlert(context.TODO(), alert2)

	// The alert should be delayed by the group_wait, so neither alert should have come through yet.
	require.NotContains(t, k.Stdout(), "bar")
	require.NotContains(t, k.Stdout(), "baz")

	// 2s is the group wait.
	time.Sleep(2 * time.Second)
	require.Contains(t, k.Stdout(), "bar")
	require.Contains(t, k.Stdout(), "baz")

	// Wait another group to make sure it doesn't fire again.
	time.Sleep(2 * time.Second)
	require.Equal(t, 1, strings.Count(k.Stdout(), "bar"))
	require.Equal(t, 1, strings.Count(k.Stdout(), "baz"))
}
