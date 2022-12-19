package raft

import (
	"context"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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
	msg := NewPostAlertsRaftLogMessage(alerts...)
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	f := r.raft.Apply(bytes, 0)
	return f.Error()
}

// GetAlerts gets all the alerts currently in the database.
func (r *RaftDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return r.db.GetAlerts(ctx)
}
