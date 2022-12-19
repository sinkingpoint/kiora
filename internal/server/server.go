package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sinkingpoint/kiora/internal/raft"
	"github.com/sinkingpoint/kiora/internal/server/apiv1"
	"github.com/sinkingpoint/kiora/internal/server/raftadmin"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

// TLSPair is a pair of paths representing the path to a certificate and private key.
type TLSPair struct {
	// CertPath is the path to the certificate.
	CertPath string

	// KeyPath is the path to the private key.
	KeyPath string
}

type serverConfig struct {
	// ListenAddress is the address for the server to listen on. Defaults to localhost:4278.
	ListenAddress string

	// ReadTimeout is the maximum amount of time the server will spend reading requests from clients. Defaults to 5 seconds.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum amount of time the server will spend writing requests to clients. Defaults to 60 seconds.
	WriteTimeout time.Duration

	// TLS is an optional pair of cert and key files that will be used to serve TLS connections.
	TLS *TLSPair
}

// NewServerConfig constructs a serverConfig with all the defaults set.
func NewServerConfig() serverConfig {
	return serverConfig{
		ListenAddress: "localhost:4278",
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  60 * time.Second,
		TLS:           nil,
	}
}

// KioraServer is a server that serves the main Kiora API.
type KioraServer struct {
	serverConfig
	db kioradb.DB
}

func NewKioraServer(conf serverConfig, db kioradb.DB) *KioraServer {
	return &KioraServer{
		serverConfig: conf,
		db:           db,
	}
}

// ListenAndServe starts the server, using TLS if set in the config. This method blocks until the server ends.
func (k *KioraServer) ListenAndServe() error {
	r := mux.NewRouter()

	apiv1.Register(r, k.db)

	if raft, ok := k.db.(*raft.RaftDB); ok {
		raftadmin.Register(r, raft.Raft)
	}

	server := http.Server{
		Addr:         k.ListenAddress,
		ReadTimeout:  k.ReadTimeout,
		WriteTimeout: k.WriteTimeout,
		Handler:      r,
	}

	var err error

	if k.TLS != nil {
		err = server.ListenAndServeTLS(k.TLS.CertPath, k.TLS.KeyPath)
	} else {
		err = server.ListenAndServe()
	}

	// ListenAndServe always returns an error, which is ErrServerClosed if cleanly exitted. Here we map
	// that into a nil for easier handling in consumers.
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
