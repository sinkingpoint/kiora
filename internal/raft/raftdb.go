package raft

import (
	"context"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/internal/server/raftadmin"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var _ kioradb.ModelWriter = &RaftDB{}

type RaftDB struct {
	myID      raft.ServerID
	raft      *raft.Raft
	transport *transport.Manager
}

func NewRaftDB(ctx context.Context, config RaftConfig, backingDB kioradb.DB) (*RaftDB, error) {
	localID := raft.ServerID(config.LocalID)
	raft, transport, err := NewRaft(ctx, config, &kioraFSM{db: backingDB})
	if err != nil {
		return nil, err
	}

	return &RaftDB{
		myID:      localID,
		raft:      raft,
		transport: transport,
	}, nil
}

func (r *RaftDB) Raft() *raft.Raft {
	return r.raft
}

// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
func (r *RaftDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	return r.applyLog(ctx, newPostAlertsRaftLogMessage(alerts...))
}

func (r *RaftDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	return r.applyLog(ctx, newPostSilencesRaftLogMessage(silences...))
}

// applyLog takes the given protobuf message, marshals it, and adds it as a log into the raft log.
func (r *RaftDB) applyLog(ctx context.Context, msg *kioraproto.RaftLogMessage) error {
	leaderAddress, leaderID := r.raft.LeaderWithID()

	if leaderID == r.myID {
		return r.applyAsLeader(ctx, msg)
	}

	return r.forwardLog(ctx, string(leaderAddress), msg)
}

// forwardLog is responsible for forwarding a log to the leader node, in the case that the node that received the log is not the leader.
func (r *RaftDB) forwardLog(ctx context.Context, leaderAddress string, msg *kioraproto.RaftLogMessage) error {
	conn, err := grpc.Dial(string(leaderAddress), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	f := r.raft.Apply(bytes, 0)
	return f.Error()
}

func (r *RaftDB) RegisterEndpoints(ctx context.Context, httpRouter *mux.Router, grpcServer *grpc.Server) error {
	r.transport.Register(grpcServer)
	raftadmin.Register(httpRouter, r.raft)
	return nil
}
