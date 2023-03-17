package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func TestFailover(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// The general flow of this test is:
	// - Spin up a cluster, with 3 nodes
	// While we still have nodes in the cluster:
	// - Send an alert into the cluster
	// - Observe that the alert gets sent
	// - Shutdown the node that sends the alert

	t.Parallel()

	alert := dummyAlert()
	resolvedAlert := dummyAlert()
	resolvedAlert.Status = model.AlertStatusResolved

	nodes := StartKioraCluster(t, 3)

	for len(nodes) > 0 {
		for i := range nodes {
			nodes[i].SendAlert(t, context.TODO(), alert)
		}

		// wait a bit for the gossip to settle.
		time.Sleep(10 * time.Second)

		foundNodeIndex := -1
		for i, n := range nodes {
			if strings.Contains(n.stdout.String(), "foo") {
				require.NoError(t, n.Stop())
				foundNodeIndex = i
				break
			}
		}

		found := foundNodeIndex >= 0
		nodeNames := []string{}
		for _, node := range nodes {
			nodeNames = append(nodeNames, node.name)
		}

		require.True(t, found, "failed to find the firing node (still have nodes: %+v)", nodeNames)
		require.NoError(t, nodes[foundNodeIndex].Stop())

		nodes = append(nodes[:foundNodeIndex], nodes[foundNodeIndex+1:]...)

		// resolve the alert so it can fire again.
		for i := range nodes {
			nodes[i].SendAlert(t, context.TODO(), resolvedAlert)
		}
	}
}
