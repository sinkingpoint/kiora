package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrouping(t *testing.T) {
	initT(t)

	k := NewKioraInstance().WithConfig(t, `digraph config {
		wait_1s [type="group_wait" duration="1s"];
		group_by_alertname [type="group_labels" labels="alertname"];
		console [type="stdout"];

		alerts -> group_by_alertname -> wait_1s -> console;
	}`).Start(t)

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

	k.SendAlert(t, context.TODO(), alert1)
	k.SendAlert(t, context.TODO(), alert2)

	// The alert should be delayed by the group_wait, so neither alert should have come through yet.
	assert.NotContains(t, k.Stdout(), "bar")
	assert.NotContains(t, k.Stdout(), "baz")

	// 2s is the group wait.
	time.Sleep(2 * time.Second)
	assert.Contains(t, k.Stdout(), "bar")
	assert.Contains(t, k.Stdout(), "baz")

	// Wait another group to make sure it doesn't fire again.
	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, strings.Count(k.Stdout(), "bar"))
	assert.Equal(t, 1, strings.Count(k.Stdout(), "baz"))
}
