package clustering

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type hasher struct{}

func (h hasher) Sum64(value []byte) uint64 {
	return xxhash.Sum64(value)
}

// Clusterer provides a way for nodes to query the cluster of nodes they belong to.
type Clusterer interface {
	AmIAuthoritativeFor(a *model.Alert) bool
}

// KioraClusterer is a Clusterer that uses a consistant hashing ring to determine Authoritative nodes.
type KioraClusterer struct {
	// the name of this node
	myName string

	// the ring that we use to work out authoritative nodes.
	ring *consistent.Consistent
}

func NewKioraClusterer(myName string) *KioraClusterer {
	return &KioraClusterer{
		myName: myName,
		ring:   nil,
	}
}

func (k *KioraClusterer) init(server consistent.Member) {
	conf := consistent.Config{
		Hasher: hasher{},
	}

	k.ring = consistent.New([]consistent.Member{
		server,
	}, conf)
}

func (k *KioraClusterer) AddServer(server Server) {
	if k.ring == nil {
		k.init(server)
		return
	}

	k.ring.Add(server)
}

func (k *KioraClusterer) RemoveServer(server Server) {
	k.ring.Remove(server.String())
}

func (k *KioraClusterer) AmIAuthoritativeFor(a *model.Alert) bool {
	member := k.ring.LocateKey(a.Labels.Bytes())

	return member != nil && member.String() == k.myName
}
