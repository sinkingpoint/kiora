package server

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"sync"

	_ "net/http/pprof"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/clustering/serf"
	"github.com/sinkingpoint/kiora/internal/pipeline"
	"github.com/sinkingpoint/kiora/internal/server/api"
	"github.com/sinkingpoint/kiora/internal/server/api/apiv1"
	"github.com/sinkingpoint/kiora/internal/server/api/promcompat"
	"github.com/sinkingpoint/kiora/internal/server/frontend"
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

	// We generate the config up here so that we have a concrete node name to pass to the clusterer.
	// TODO(cdouch): This generates a random node name. Allow the user to specify a node name.
	config := serf.DefaultConfig()

	ringClusterer := clustering.NewRingClusterer(config.NodeName, clusterAddress)
	if len(conf.ClusterShardLabels) > 0 {
		ringClusterer.SetShardLabels(conf.ClusterShardLabels)
	}

	delegate := pipeline.NewDBEventDelegate(db)
	config.EventDelegate = delegate
	config.ListenURL = conf.ClusterListenAddress
	config.BootstrapPeers = conf.BootstrapPeers
	config.ClustererDelegate = ringClusterer
	config.Logger = conf.Logger
	config.DBDelegate = serf.NewDBDelegate(db, conf.Logger)

	broadcaster, err := serf.NewSerfBroadcaster(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct broadcaster")
	}

	bus := services.NewKioraBus(db, broadcaster, config.Logger, conf.ServiceConfig)

	services := services.NewBackgroundServices()
	services.RegisterService(broadcaster)
	services.RegisterService(notify.NewNotifyService(notify_config.NewClusterNotifier(ringClusterer, conf.ServiceConfig), bus))
	services.RegisterService(timeout.NewTimeoutService(bus))
	services.RegisterService(delegate)

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
		k.backgroundServices.Shutdown(context.Background())
	})
}

// ListenAndServe starts the server, using TLS if set in the config. This method blocks until the server ends.
func (k *KioraServer) ListenAndServe() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := k.listenAndServeHTTP(context.Background()); err != nil {
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

		if k.httpServer != nil {
			k.httpServer.Shutdown(context.Background()) //nolint:errcheck
		}

		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (k *KioraServer) listenAndServeHTTP(ctx context.Context) error {
	router := mux.NewRouter()
	router.PathPrefix("/debug/").Handler(http.DefaultServeMux)

	api := api.NewAPIImpl(k.bus, k.clusterer)
	apiv1.Register(router, api, k.serverConfig.Logger)
	promcompat.Register(router, api, k.serverConfig.Logger)

	frontend.Register(router)

	runtime.SetMutexProfileFraction(5)

	k.httpServer = &http.Server{
		Addr:         k.HTTPListenAddress,
		ReadTimeout:  k.ReadTimeout,
		WriteTimeout: k.WriteTimeout,
		Handler:      handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router),
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
