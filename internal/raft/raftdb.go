package raft

import (
	"context"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/internal/server/raftadmin"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var _ kioradb.ModelWriter = &RaftDB{}

type RaftDB struct {
	myID         raft.ServerID
	raft         *raft.Raft
	transport    *transport.Manager
	dispatchChan chan *kioraproto.RaftLogMessage
}

func NewRaftDB(ctx context.Context, config RaftConfig, backingDB kioradb.DB) (*RaftDB, error) {
	localID := raft.ServerID(config.LocalID)
	raft, transport, err := NewRaft(ctx, config, &kioraFSM{db: backingDB})
	if err != nil {
		return nil, err
	}

	db := RaftDB{
		myID:         localID,
		raft:         raft,
		transport:    transport,
		dispatchChan: make(chan *kioraproto.RaftLogMessage, 500), // TODO(cdouch): This capacity is arbitrary. Should benchmark it.
	}

	go func() {
		for msg := range db.dispatchChan {
			if err := db.applyLog(context.Background(), msg); err != nil {
				log.Err(err).Msg("failed to apply log")
			}
		}
	}()

	return &db, nil
}

func (r *RaftDB) Raft() *raft.Raft {
	return r.raft
}

// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
func (r *RaftDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	r.dispatchChan <- newPostAlertsRaftLogMessage(alerts...)
	return nil
}

func (r *RaftDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	r.dispatchChan <- newPostSilencesRaftLogMessage(silences...)
	return nil
}

// applyLog takes the given protobuf message, marshals it, and adds it as a log into the raft log.
func (r *RaftDB) applyLog(ctx context.Context, msg *kioraproto.RaftLogMessage) error {
	ctx, span := tracing.Tracer().Start(ctx, "RaftDB.applyLog")
	defer span.End()

	leaderAddress, leaderID := r.raft.LeaderWithID()

	if leaderID == r.myID {
		return r.applyAsLeader(ctx, msg)
	}

	span.AddEvent("forwarding")

	return r.forwardLog(ctx, string(leaderAddress), msg)
}

// forwardLog is responsible for forwarding a log to the leader node, in the case that the node that received the log is not the leader.
func (r *RaftDB) forwardLog(ctx context.Context, leaderAddress string, msg *kioraproto.RaftLogMessage) error {
	conn, err := grpc.Dial(string(leaderAddress), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if err != nil {
		return err
	}
	defer conn.Close()

	msg.From = string(r.myID)
	client := kioraproto.NewRaftApplierClient(conn)
	_, err = client.ApplyLog(ctx, msg)

	return err
}

// applyAsLeader gets called to apply a log when this node is the leader of the cluster. When inside this method
// it can be assumed that this node is the leader, and thus methods that must be called on the raft leader are safe.
func (r *RaftDB) applyAsLeader(ctx context.Context, msg *kioraproto.RaftLogMessage) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return r.raft.ApplyLogCtx(ctx, raft.Log{Data: bytes}).Error()
}

func (r *RaftDB) RegisterEndpoints(ctx context.Context, httpRouter *mux.Router, grpcServer *grpc.Server) error {
	r.transport.Register(grpcServer)
	raftadmin.Register(httpRouter, r.raft)
	return nil
}
