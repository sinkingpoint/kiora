package server

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/internal/raft"
	"github.com/sinkingpoint/kiora/internal/server/apiv1"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	GRPCListenAddress string

	// ReadTimeout is the maximum amount of time the server will spend reading requests from clients. Defaults to 5 seconds.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum amount of time the server will spend writing requests to clients. Defaults to 60 seconds.
	WriteTimeout time.Duration

	// TLS is an optional pair of cert and key files that will be used to serve TLS connections.
	TLS *TLSPair

	RaftConfig raft.RaftConfig
}

// NewServerConfig constructs a serverConfig with all the defaults set.
func NewServerConfig() serverConfig {
	return serverConfig{
		HTTPListenAddress: "localhost:4278",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      60 * time.Second,
		TLS:               nil,
		RaftConfig:        raft.DefaultRaftConfig(),
	}
}

// KioraServer is a server that serves the main Kiora API.
type KioraServer struct {
	kioraproto.UnimplementedKioraServer
	serverConfig

	httpServer *http.Server
	grpcServer *grpc.Server

	broadcaster clustering.Broadcaster
	db          kioradb.DB

	shutdownOnce sync.Once
}

func NewKioraServer(conf serverConfig, db kioradb.DB) (*KioraServer, error) {
	broadcaster, err := raft.NewRaftBroadcaster(context.Background(), conf.RaftConfig, db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise raft")
	}

	return &KioraServer{
		db:           db,
		serverConfig: conf,
		broadcaster:  broadcaster,
	}, nil
}

func (k *KioraServer) Shutdown() {
	k.shutdownOnce.Do(func() {
		// todo: add more synonyms of stop.
		k.httpServer.Shutdown(context.Background()) //nolint:errcheck
		k.grpcServer.GracefulStop()
	})
}

// ListenAndServe starts the server, using TLS if set in the config. This method blocks until the server ends.
func (k *KioraServer) ListenAndServe() error {
	k.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))

	httpRouter := mux.NewRouter()

	if err := k.broadcaster.RegisterEndpoints(context.Background(), httpRouter, k.grpcServer); err != nil {
		return errors.Wrap(err, "failed to register broadcaster endpoints")
	}

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
		if err := k.listenAndServeGRPC(); err != nil {
			log.Err(err).Msg("Error shutting down GRPC server")
		}

		wg.Done()

		log.Info().Msg("GRPC Server Shut Down")
	}()

	wg.Wait()

	return nil
}

func (k *KioraServer) listenAndServeGRPC() error {
	listener, err := net.Listen("tcp", k.serverConfig.GRPCListenAddress)
	if err != nil {
		return err
	}

	kioraproto.RegisterKioraServer(k.grpcServer, k)
	reflection.Register(k.grpcServer)
	return k.grpcServer.Serve(listener)
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

// ApplyLog is the handler that Log Messages that have been received on a follower and forwarded to the Raft leader for applying.
func (k *KioraServer) ApplyLog(ctx context.Context, log *kioraproto.KioraLogMessage) (*kioraproto.KioraLogReply, error) {
	switch msg := log.Log.(type) {
	case *kioraproto.KioraLogMessage_Alerts:
		alerts := []model.Alert{}
		for _, protoAlert := range msg.Alerts.Alerts {
			alert := model.Alert{}

			if err := alert.DeserializeFromProto(protoAlert); err != nil {
				return nil, err
			}

			alerts = append(alerts, alert)
		}

		if err := k.broadcaster.BroadcastAlerts(ctx, alerts...); err != nil {
			return nil, err
		}
	}

	return &kioraproto.KioraLogReply{}, nil
}

func (k *KioraServer) Heartbeat(ctx context.Context, hearbeat *kioraproto.HeartbeatMessage) (*kioraproto.HeartbeatReply, error) {
	return &kioraproto.HeartbeatReply{}, nil
}
