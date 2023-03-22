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

// Test that a cluster of three nodes, when an alert is posted to one, that alert gets posted to all.
func TestKioraClusterAlerts(t *testing.T) {
	initT(t)

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
	initT(t)

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

// Tests that we can fire alert, acknowledge it, and get the acknowledgement back.
func TestAcknowledgementGetsRegistered(t *testing.T) {
	initT(t)

	// Send an alert.
	alert := dummyAlert()
	nodes := StartKioraCluster(t, 3)
	nodes[0].SendAlert(t, context.TODO(), alert)

	time.Sleep(2 * time.Second)

	// Make sure the alert exists.
	alerts := nodes[0].GetAlerts(t, context.TODO())
	require.Len(t, alerts, 1)

	// Send an acknowledgement for that alert.
	nodes[0].SendAlertAcknowledgement(t, context.TODO(), ackRequest{
		AlertAcknowledgement: model.AlertAcknowledgement{
			Creator: "test_creator",
			Comment: "test_comment",
		},
		AlertID: alerts[0].ID,
	})

	time.Sleep(2 * time.Second)

	// Get the alerts again and make sure our acknowledgement is there.
	alerts = nodes[0].GetAlerts(t, context.TODO())
	require.Len(t, alerts, 1)

	assert.NotNil(t, alerts[0].Acknowledgement)
	assert.Equal(t, "test_creator", alerts[0].Acknowledgement.Creator)
	assert.Equal(t, "test_comment", alerts[0].Acknowledgement.Comment)
	assert.Equal(t, model.AlertStatusAcked, alerts[0].Status)
}

// Test that we can add a silence, and that it prevents the alert from being sent.
func TestSilencesSilence(t *testing.T) {
	initT(t)

	silencedAlert := dummyAlert()
	nonSilencedAlert := dummyAlert()
	silence := dummySilence()

	silencedAlert.Labels["silenced"] = "true"
	nonSilencedAlert.Labels["silenced"] = "false"
	silence.Matchers = append(silence.Matchers, model.Matcher{
		Label: "silenced",
		Value: "true",
	})

	nodes := StartKioraCluster(t, 3)

	// Send a silence.
	nodes[0].SendSilence(t, context.TODO(), silence)
	time.Sleep(2 * time.Second)

	// Send in an alert that shouldn't be silenced.
	nodes[0].SendAlert(t, context.TODO(), nonSilencedAlert)
	// Send an alert that should be silenced.
	nodes[0].SendAlert(t, context.TODO(), silencedAlert)
	time.Sleep(2 * time.Second)

	// Make sure the alert exists.
	alerts := nodes[0].GetAlerts(t, context.TODO())
	require.Len(t, alerts, 2)

	foundSilenced := 0
	for _, alert := range alerts {
		if alert.Status == model.AlertStatusSilenced {
			foundSilenced += 1
			assert.Equal(t, "true", alert.Labels["silenced"])
		} else {
			assert.Equal(t, "false", alert.Labels["silenced"])
		}
	}

	assert.Equal(t, 1, foundSilenced, "Expected one silenced alert")

	foundAlerted := false
	// Assert that no node sent the alert for the silenced alert.
	for _, node := range nodes {
		assert.NotContains(t, node.Stdout(), `"silenced":"true"`)

		if strings.Contains(node.Stdout(), `"silenced":"false"`) {
			if foundAlerted {
				t.Fatal("Expected only one node to send the alert for the non-silenced alert")
			}
			foundAlerted = true
		}
	}

	// Assert that one node sent the alert for the silenced alert.
	assert.True(t, foundAlerted, "Expected one node to send the alert for the non-silenced alert")
}

// Test that if we send an alert, and then silence it, the alert status gets updated.
func TestSilencesSilenceAfterAlert(t *testing.T) {
	initT(t)
	alert := dummyAlert()
	silence := dummySilence()

	nodes := StartKioraCluster(t, 3)
	nodes[0].SendAlert(t, context.TODO(), alert)
	time.Sleep(2 * time.Second)

	// Send a silence.
	nodes[0].SendSilence(t, context.TODO(), silence)
	time.Sleep(2 * time.Second)

	// Make sure the alert exists and is silenced.
	alerts := nodes[0].GetAlerts(t, context.TODO())
	require.Len(t, alerts, 1)
	assert.Equal(t, model.AlertStatusSilenced, alerts[0].Status)
}
