package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/internal/kiora"
	"github.com/sinkingpoint/kiora/internal/raft"
	"github.com/sinkingpoint/kiora/internal/server/apiv1"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
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

	NotifyConfig kiora.NotifierConfig

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

// assembleProcessor is responsible for constructing the KioraProcessor that pre-processes models before they enter the main flow.
func assemblePreProcessor(conf *serverConfig, broadcaster kioradb.ModelWriter, db kioradb.DB) *kiora.KioraProcessor {
	processor := kiora.NewKioraProcessor(db, broadcaster)

	// For now, just broadcast everything that comes in.
	broadcastProcessor := kiora.BroadcastProcessor{}
	processor.AddAlertProcessor(&broadcastProcessor)
	processor.AddSilenceProccessor(&broadcastProcessor)

	return processor
}

// assemblePostProcessor is responsible for constructing the KioraProcessor that processes models _after_ they have been broadcasted.
func assemblePostProcessor(conf *serverConfig, broadcaster kioradb.ModelWriter, db kioradb.DB) *kiora.KioraProcessor {
	processor := kiora.NewKioraProcessor(db, broadcaster)

	localForwarder := kiora.LocalForwarderProcessor{}
	processor.AddAlertProcessor(kiora.NewSilenceApplier())
	processor.AddAlertProcessor(kiora.NewNotifierProcessor(conf.RaftConfig.LocalID, conf.NotifyConfig))
	processor.AddAlertProcessor(&localForwarder)
	processor.AddSilenceProccessor(&localForwarder)

	return processor
}

// KioraServer is a server that serves the main Kiora API.
type KioraServer struct {
	kioraproto.UnimplementedRaftApplierServer
	serverConfig
	broadcaster *raft.RaftDB
	db          *kiora.KioraProcessor
}

func NewKioraServer(conf serverConfig, db kioradb.DB) (*KioraServer, error) {
	postProcessor := assemblePostProcessor(&conf, nil, db)
	broadcaster, err := raft.NewRaftDB(context.Background(), conf.RaftConfig, postProcessor)
	if err != nil {
		return nil, err
	}

	// this makes a weird loop where we could go broadcast -> postprocessor -> broadcast infinitely.
	// TODO(cdouch): detangle this if the circular dependency proves unweildy.
	postProcessor.Broadcast = broadcaster

	return &KioraServer{
		serverConfig: conf,
		db:           assemblePreProcessor(&conf, broadcaster, db),
		broadcaster:  broadcaster,
	}, nil
}

// ListenAndServe starts the server, using TLS if set in the config. This method blocks until the server ends.
func (k *KioraServer) ListenAndServe() error {
	errChan := make(chan error)

	grpcServer := grpc.NewServer()
	httpRouter := mux.NewRouter()

	if err := k.db.RegisterEndpoints(context.Background(), httpRouter, grpcServer); err != nil {
		return err
	}

	if err := k.broadcaster.RegisterEndpoints(context.Background(), httpRouter, grpcServer); err != nil {
		return err
	}

	go func() {
		errChan <- k.listenAndServeHTTP(httpRouter)
	}()

	go func() {
		errChan <- k.listenAndServeGRPC(grpcServer)
	}()

	return <-errChan
}

func (k *KioraServer) listenAndServeGRPC(server *grpc.Server) error {
	listener, err := net.Listen("tcp", k.serverConfig.GRPCListenAddress)
	if err != nil {
		return err
	}

	kioraproto.RegisterRaftApplierServer(server, k)
	reflection.Register(server)
	return server.Serve(listener)
}

func (k *KioraServer) listenAndServeHTTP(r *mux.Router) error {
	apiv1.Register(r, k.db)

	httpServer := http.Server{
		Addr:         k.HTTPListenAddress,
		ReadTimeout:  k.ReadTimeout,
		WriteTimeout: k.WriteTimeout,
		Handler:      r,
	}

	var err error

	if k.TLS != nil {
		err = httpServer.ListenAndServeTLS(k.TLS.CertPath, k.TLS.KeyPath)
	} else {
		err = httpServer.ListenAndServe()
	}

	// ListenAndServe always returns an error, which is ErrServerClosed if cleanly exitted. Here we map
	// that into a nil for easier handling in consumers.
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// ApplyLog is the handler that processes forwarded RaftLogs from follower nodes to the leader. When this function is called,
// it can be assumed that the current node is the leader of the raft cluster.
func (k *KioraServer) ApplyLog(ctx context.Context, log *kioraproto.RaftLogMessage) (*kioraproto.RaftLogReply, error) {
	switch msg := log.Log.(type) {
	case *kioraproto.RaftLogMessage_Alerts:
		modelAlerts := make([]model.Alert, 0, len(msg.Alerts.Alerts))
		for _, protoAlert := range msg.Alerts.Alerts {
			alert := model.Alert{
				AuthNode: log.From,
			}

			if err := alert.DeserializeFromProto(protoAlert); err != nil {
				return nil, err
			}

			modelAlerts = append(modelAlerts, alert)
		}

		return &kioraproto.RaftLogReply{}, k.db.ProcessAlerts(ctx, modelAlerts...)
	case *kioraproto.RaftLogMessage_Silences:
	}

	return &kioraproto.RaftLogReply{}, nil
}
