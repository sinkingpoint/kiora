package clustering

import (
	"context"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type hasher struct{}

func (h hasher) Sum64(bytes []byte) uint64 {
	return xxhash.Sum64(bytes)
}

type kioraMember struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (k *kioraMember) String() string {
	return k.Name
}

var (
	_ = Clusterer(&RingClusterer{})
	_ = ClustererDelegate(&RingClusterer{})
)

// RingClusterer is a clusterer that keeps track of nodes in a consistent hash ring.
type RingClusterer struct {
	me        consistent.Member
	ring      *consistent.Consistent
	shardKeys []string
}

// NewRingClusterer constructs a new RingClusterer, with the given name and address.
// This name and address _must_ be the same as the node in the underlying Cluster in order to properly shard alerts.
func NewRingClusterer(myName, myAddress string) *RingClusterer {
	me := &kioraMember{
		Name:    myName,
		Address: myAddress,
	}

	config := consistent.Config{
		Hasher: hasher{},
	}

	return &RingClusterer{
		me:   me,
		ring: consistent.New([]consistent.Member{me}, config),
	}
}

func (r *RingClusterer) SetShardLabels(keys []string) {
	r.shardKeys = keys
}

func (r *RingClusterer) IsAuthoritativeFor(ctx context.Context, a *model.Alert) bool {
	return r.GetAuthoritativeNode(ctx, a) == r.me
}

// GetAuthoritativeNode returns the node that is authoritative for the given alert.
func (r *RingClusterer) GetAuthoritativeNode(ctx context.Context, a *model.Alert) consistent.Member {
	if len(r.shardKeys) == 0 {
		return r.ring.LocateKey(a.Labels.Bytes())
	}

	labels := a.Labels.Subset(r.shardKeys...)
	return r.ring.LocateKey(labels.Bytes())
}

func (r *RingClusterer) AddNode(name, address string) {
	r.ring.Add(&kioraMember{
		Name:    name,
		Address: address,
	})
}

func (r *RingClusterer) RemoveNode(name string) {
	r.ring.Remove(name)
}

func (r *RingClusterer) Nodes() []any {
	members := r.ring.GetMembers()
	nodes := make([]any, 0, len(members))

	for _, node := range members {
		nodes = append(nodes, node)
	}

	return nodes
}
