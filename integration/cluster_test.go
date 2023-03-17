package integration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dummyAlert() model.Alert {
	return model.Alert{
		Labels: model.Labels{
			"foo": "bar",
		},
		Annotations: map[string]string{},
		Status:      model.AlertStatusFiring,
		StartTime:   time.Now(),
	}
}

// Test that a cluster of three nodes, when an alert is posted to one, that alert gets posted to all.
func TestKioraClusterAlerts(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	nodes := StartKioraCluster(t, 3)

	nodes[0].SendAlert(t, context.TODO(), dummyAlert())

	// Wait a bit for the gossip to settle.
	time.Sleep(1 * time.Second)

	for _, k := range nodes {
		reqURL := k.GetHTTPURL("/api/v1/alerts")
		resp, err := http.Get(reqURL)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Contains(t, string(respBytes), "foo")
	}
}

func TestClusterAlertOnlySentOnce(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()

	alert := dummyAlert()

	nodes := StartKioraCluster(t, 3)
	for i := range nodes {
		nodes[i].SendAlert(t, context.TODO(), alert)
	}

	// Wait a bit for the gossip to settle.
	time.Sleep(1 * time.Second)

	found := 0
	notFound := 0

	for i := range nodes {
		if strings.Contains(nodes[i].stdout.String(), "foo") {
			t.Logf("Node %d notified for alert", i)
			found += 1
		} else {
			t.Logf("Node %d did not for alert", i)
			notFound += 1
		}
	}

	assert.Equal(t, 1, found, "Expected only one notification")
	assert.Equal(t, len(nodes)-1, notFound, "Excepted two nodes to not send the notification")

	found = 0
	notFound = 0
	totalFound := 0

	// Apply the alert _a second time_ to test deduplication logic.
	for i := range nodes {
		nodes[i].SendAlert(t, context.TODO(), alert)
	}

	// Wait for the gossip to settle.
	time.Sleep(1 * time.Second)

	for i := range nodes {
		if strings.Contains(nodes[i].Stdout(), "foo") {
			t.Logf("Node %d notified for alert", i)
			found += 1
			totalFound += strings.Count(nodes[i].Stdout(), "foo")
		} else {
			t.Logf("Node %d did not for alert", i)
			notFound += 1
		}
	}

	assert.Equal(t, 1, found, "Expected only one notification")
	assert.Equal(t, 1, totalFound, "Expected only one notification")
	assert.Equal(t, len(nodes)-1, notFound, "Excepted two nodes to not send the notification")
}
