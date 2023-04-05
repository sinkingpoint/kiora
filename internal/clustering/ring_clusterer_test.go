package clustering_test

import (
	"context"
	"testing"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
)

func TestRingClustererSharding(t *testing.T) {
	clusterer := clustering.NewRingClusterer("a", "a")
	// Add a bunch of nodes to decrease the likelyhood that we get proper sharding just by chance.
	for i := 'b'; i < 'z'; i++ {
		clusterer.AddNode(string(i), string(i))
	}

	for i := 'A'; i < 'Z'; i++ {
		clusterer.AddNode(string(i), string(i))
	}

	clusterer.SetShardLabels([]string{"foo"})
	authA := clusterer.GetAuthoritativeNode(context.TODO(), &model.Alert{
		Labels: model.Labels{
			"foo": "bar",
			"bar": "baz",
		},
	})

	authB := clusterer.GetAuthoritativeNode(context.TODO(), &model.Alert{
		Labels: model.Labels{
			"foo": "bar",
			"bar": "foo",
		},
	})

	authC := clusterer.GetAuthoritativeNode(context.TODO(), &model.Alert{
		Labels: model.Labels{
			"foo": "baz",
			"bar": "baz",
		},
	})

	assert.Equal(t, authA, authB)
	assert.NotEqual(t, authA, authC)
}
