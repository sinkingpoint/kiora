package server

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

// TLSPair is a pair of paths representing the path to a certificate and private key.
type TLSPair struct {
	// CertPath is the path to the certificate.
	CertPath string

	// KeyPath is the path to the private key.
	KeyPath string
}

type serverConfig struct {
	// HTTPListenAddress is the address for the server to listen on. Defaults to localhost:4278.
	HTTPListenAddress string

	ClusterListenAddress string

	// ClusterShardLabels is the set of labels that will be used to determine which node in a cluster will send a given alert.
	// Defaults to an empty list which will shard by every label, effectively causing a random assignment across the cluster. Setting this can improve alert grouping,
	// at the cost of a potentially unbalanced cluster.
	ClusterShardLabels []string

	// BootstrapPeers is the set of peers to bootstrap the cluster with. Defaults to an empty list which means that the node will not join a cluster.
	BootstrapPeers []string

	// ServiceConfig is the config that will determine how data flows through the kiora instance.
	ServiceConfig config.Config

	// ReadTimeout is the maximum amount of time the server will spend reading requests from clients. Defaults to 5 seconds.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum amount of time the server will spend writing requests to clients. Defaults to 60 seconds.
	WriteTimeout time.Duration

	// TLS is an optional pair of cert and key files that will be used to serve TLS connections.
	TLS *TLSPair

	Logger zerolog.Logger
}

// NewServerConfig constructs a serverConfig with all the defaults set.
func NewServerConfig() serverConfig {
	return serverConfig{
		HTTPListenAddress: "localhost:4278",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      60 * time.Second,
		TLS:               nil,
	}
}
