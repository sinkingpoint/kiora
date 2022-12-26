package raft

import (
	"context"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ kioradb.DB = &RaftDB{}

type RaftDB struct {
	raft      *raft.Raft
	db        kioradb.DB
	transport *transport.Manager
}

func NewRaftDB(ctx context.Context, config raftConfig, backingDB kioradb.DB) (*RaftDB, error) {
	raft, transport, err := NewRaft(ctx, config, &kioraFSM{db: backingDB})
	if err != nil {
		return nil, err
	}

	return &RaftDB{
		raft:      raft,
		transport: transport,
		db:        backingDB,
	}, nil
}

func (r *RaftDB) RegisterGRPC(s *grpc.Server) {
	r.transport.Register(s)
}

func (r *RaftDB) Raft() *raft.Raft {
	return r.raft
}

// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
func (r *RaftDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	return r.applyLog(newPostAlertsRaftLogMessage(alerts...))
}

func (r *RaftDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	return r.applyLog(newPostSilencesRaftLogMessage(silences...))
}

// GetAlerts gets all the alerts currently in the database.
func (r *RaftDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return r.db.GetAlerts(ctx)
}

func (r *RaftDB) GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error) {
	return r.db.GetExistingAlert(ctx, labels)
}

// applyLog takes the given protobuf message, marshals it, and adds it as a log into the raft log.
func (r *RaftDB) applyLog(msg protoreflect.ProtoMessage) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	f := r.raft.Apply(bytes, 0)
	return f.Error()
}
