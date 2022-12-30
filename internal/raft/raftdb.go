package raft

import (
	"context"
	"fmt"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var _ kioradb.DB = &RaftDB{}

type RaftDB struct {
	myID      raft.ServerID
	raft      *raft.Raft
	db        kioradb.DB
	transport *transport.Manager
}

func NewRaftDB(ctx context.Context, config raftConfig, backingDB kioradb.DB) (*RaftDB, error) {
	localID := raft.ServerID(config.LocalID)
	raft, transport, err := NewRaft(ctx, config, &kioraFSM{db: backingDB})
	if err != nil {
		return nil, err
	}

	return &RaftDB{
		myID:      localID,
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
	return r.applyLog(ctx, newPostAlertsRaftLogMessage(alerts...))
}

func (r *RaftDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	return r.applyLog(ctx, newPostSilencesRaftLogMessage(silences...))
}

// GetAlerts gets all the alerts currently in the database.
func (r *RaftDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return r.db.GetAlerts(ctx)
}

func (r *RaftDB) GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error) {
	return r.db.GetExistingAlert(ctx, labels)
}

func (r *RaftDB) GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error) {
	return r.db.GetSilences(ctx, labels)
}

// applyLog takes the given protobuf message, marshals it, and adds it as a log into the raft log.
func (r *RaftDB) applyLog(ctx context.Context, msg *kioraproto.RaftLogMessage) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	leaderAddress, leaderID := r.raft.LeaderWithID()

	if leaderID == r.myID {
		f := r.raft.Apply(bytes, 0)
		return f.Error()
	}

	conn, err := grpc.Dial(string(leaderAddress), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := kioraproto.NewRaftApplierClient(conn)
	_, err = client.ApplyLog(ctx, msg)

	if err != nil {
		panic(fmt.Sprintf("FINDME: %q %q\n", ctx, msg))
	}

	return nil
}
