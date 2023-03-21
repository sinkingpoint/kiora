package server

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"sync"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/clustering/serf"
	"github.com/sinkingpoint/kiora/internal/kiora/pipeline"
	"github.com/sinkingpoint/kiora/internal/server/apiv1"
	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/internal/services/notify"
	"github.com/sinkingpoint/kiora/internal/services/notify/notify_config"
	"github.com/sinkingpoint/kiora/internal/services/timeout"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

// KioraServer is a server that serves the main Kiora API.
type KioraServer struct {
	serverConfig

	httpServer *http.Server

	bus       services.Bus
	clusterer clustering.Clusterer

	backgroundServices *services.BackgroundServices

	shutdownOnce sync.Once
}

func NewKioraServer(conf serverConfig, db kioradb.DB) (*KioraServer, error) {
	// Resolve the cluster address into an actual address so that it's consistent with other members of the cluster
	// that communicate over an actual resolved address.
	clusterAddress, err := resolveConcreteAddress(conf.ClusterListenAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode cluster listen address")
	}

	config := serf.DefaultConfig()

	ringClusterer := clustering.NewRingClusterer(config.NodeName, clusterAddress)

	config.ListenURL = conf.ClusterListenAddress
	config.BootstrapPeers = conf.BootstrapPeers
	config.ClustererDelegate = ringClusterer
	config.EventDelegate = pipeline.NewDBEventDelegate(db)
	broadcaster, err := serf.NewSerfBroadcaster(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct broadcaster")
	}

	bus := services.NewKioraBus(db, broadcaster, conf.ServiceConfig)

	services := services.NewBackgroundServices()
	services.RegisterService(broadcaster)
	services.RegisterService(notify.NewNotifyService(notify_config.NewClusterNotifier(ringClusterer, conf.ServiceConfig), bus))
	services.RegisterService(timeout.NewTimeoutService(bus))

	return &KioraServer{
		serverConfig:       conf,
		bus:                bus,
		clusterer:          ringClusterer,
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
	apiv1.Register(r, k.bus, k.clusterer)

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

// resolveConcreteAddress takes a listen address like `localhost:4278` and resolves
// it into a concrete address like `[::]:4278`
func resolveConcreteAddress(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", err
	}

	localIP, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(localIP.String(), port), nil
}
