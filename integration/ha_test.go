package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailover(t *testing.T) {
	// The general flow of this test is:
	// - Spin up a cluster, with 3 nodes
	// While we still have nodes in the cluster:
	// - Send an alert into the cluster
	// - Observe that the alert gets sent
	// - Shutdown the node that sends the alert

	t.SkipNow()

	t.Parallel()

	alert := dummyAlert()

	nodes := StartKioraCluster(t, 3)

	for len(nodes) > 0 {
		for i := range nodes {
			nodes[i].SendAlert(t, context.TODO(), alert)
		}

		// Wait a bit for raft to converge.
		time.Sleep(1 * time.Second)

		foundAndShutdown := -1
		for i, n := range nodes {
			if strings.Contains(n.stdout.String(), "foo") {
				require.NoError(t, n.Stop())
				foundAndShutdown = i
				break
			}
		}

		assert.Greater(t, foundAndShutdown, -1, "failed to find the firing node (still have nodes: %+v)", nodes)

		nodes = append(nodes[:foundAndShutdown], nodes[foundAndShutdown+1:]...)
	}
}
