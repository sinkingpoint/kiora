package server

import (
	"context"
	"net/http"
	"runtime"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/clustering/serf"
	"github.com/sinkingpoint/kiora/internal/server/apiv1"
	"github.com/sinkingpoint/kiora/internal/server/services"
	"github.com/sinkingpoint/kiora/internal/services/notify"
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
	// HTTPListenAddress is the address for the server to listen on. Defaults to localhost:4278.
	HTTPListenAddress string

	ClusterListenAddress string
	BootstrapPeers       []string
	NotifierConfig       notify.NotifierConfig

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
		HTTPListenAddress: "localhost:4278",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      60 * time.Second,
		TLS:               nil,
	}
}

// KioraServer is a server that serves the main Kiora API.
type KioraServer struct {
	serverConfig

	httpServer *http.Server

	broadcaster clustering.Broadcaster
	db          kioradb.DB

	backgroundServices *services.BackgroundServices

	shutdownOnce sync.Once
}

func NewKioraServer(conf serverConfig, db kioradb.DB) (*KioraServer, error) {
	config := serf.DefaultConfig()
	ringClusterer := clustering.NewRingClusterer(config.NodeName, "")

	config.ListenURL = conf.ClusterListenAddress
	config.BootstrapPeers = conf.BootstrapPeers
	config.ClustererDelegate = ringClusterer
	broadcaster, err := serf.NewSerfBroadcaster(config, db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct broadcaster")
	}

	services := services.NewBackgroundServices()
	services.RegisterService(broadcaster)
	services.RegisterService(notify.NewNotifyService(notify.NewClusterNotifier(ringClusterer, conf.NotifierConfig), db, broadcaster))

	return &KioraServer{
		db:                 db,
		serverConfig:       conf,
		broadcaster:        broadcaster,
		backgroundServices: services,
	}, nil
}

func (k *KioraServer) Shutdown() {
	k.shutdownOnce.Do(func() {
		// todo: add more synonyms of stop.
		if k.httpServer != nil {
			k.httpServer.Shutdown(context.Background()) //nolint:errcheck
		}

		k.backgroundServices.Shutdown(context.Background())
	})
}

// ListenAndServe starts the server, using TLS if set in the config. This method blocks until the server ends.
func (k *KioraServer) ListenAndServe() error {
	httpRouter := mux.NewRouter()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := k.listenAndServeHTTP(httpRouter); err != nil {
			log.Err(err).Msg("Error shutting down HTTP server")
		}

		wg.Done()

		log.Info().Msg("HTTP Server Shut Down")
	}()

	go func() {
		if err := k.backgroundServices.Run(context.Background()); err != nil {
			log.Err(err).Msg("background services failed")
			k.Shutdown()
		}

		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (k *KioraServer) listenAndServeHTTP(r *mux.Router) error {
	apiv1.Register(r, k.db, k.broadcaster)

	runtime.SetMutexProfileFraction(5)
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	k.httpServer = &http.Server{
		Addr:         k.HTTPListenAddress,
		ReadTimeout:  k.ReadTimeout,
		WriteTimeout: k.WriteTimeout,
		Handler:      r,
	}

	var err error

	if k.TLS != nil {
		err = k.httpServer.ListenAndServeTLS(k.TLS.CertPath, k.TLS.KeyPath)
	} else {
		err = k.httpServer.ListenAndServe()
	}

	// ListenAndServe always returns an error, which is ErrServerClosed if cleanly exitted. Here we map
	// that into a nil for easier handling in consumers.
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
