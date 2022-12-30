package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that Kiora doesn't immediatly exit.
func TestKioraStart(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()
	kiora := NewKioraInstance()
	require.NoError(t, kiora.Start(t))
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	assert.Equal(t, context.DeadlineExceeded, kiora.WaitForExit(ctx), "StdErr: %q", kiora.Stderr())
}

// Test that a post to a Kiora instance stores the alert.
func TestKioraAlertPost(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kiora := NewKioraInstance("--raft.bootstrap")
	require.NoError(t, kiora.Start(t))
	require.NoError(t, kiora.WaitUntilLeader(ctx))

	referenceTime, err := time.Parse(time.RFC3339, "2022-12-13T21:55:12Z")
	require.NoError(t, err)

	resp, err := http.Post(kiora.GetURL("/api/v1/alerts"), "application/json", strings.NewReader(fmt.Sprintf(`[{
	"labels": {},
	"annotations": {},
	"status": "firing",
	"startTime": "%s"
}]`, referenceTime.Format(time.RFC3339))))
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusAccepted, resp.StatusCode, string(body))
}

// Test that a cluster of three nodes comes up properly, with a leader.
func TestKioraCluster(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()

	nodes := StartKioraCluster(t, 3)
	for _, k := range nodes {
		status := k.ClusterStatus(t)
		assert.Contains(t, status, `"is_leader":true`)
		assert.Equal(t, 3, strings.Count(status, `"id"`), status)
	}
}

// Test that a cluster of three nodes, when an alert is posted to one, that alert gets posted to all.
func TestKioraClusterAlerts(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()

	nodes := StartKioraCluster(t, 3)
	requestURL := nodes[0].GetURL("/api/v1/alerts")

	referenceTime, err := time.Parse(time.RFC3339, "2022-12-13T21:55:12Z")
	require.NoError(t, err)

	resp, err := http.Post(requestURL, "application/json", strings.NewReader(fmt.Sprintf(`[{
	"labels": {
		"foo":"bar"
	},
	"annotations": {},
	"status": "firing",
	"startTime": "%s"
}]`, referenceTime.Format(time.RFC3339))))

	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	t.Logf("Response: %q", string(body))
	require.Equal(t, http.StatusAccepted, resp.StatusCode)

	// Wait for raft to apply the log.
	time.Sleep(1 * time.Second)

	for _, k := range nodes {
		reqURL := k.GetURL("/api/v1/alerts")
		resp, err := http.Get(reqURL)
		require.NoError(t, err)

		respBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Contains(t, string(respBytes), "foo")
	}
}
