package raft

import (
	"context"

	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/protobuf/proto"
)

var _ kioradb.DB = &RaftDB{}

type RaftDB struct {
	raft *raft.Raft
	db   kioradb.DB
}

func NewRaftDB(ctx context.Context, config raftConfig, backingDB kioradb.DB) (*RaftDB, error) {
	raft, err := NewRaft(ctx, config, &alertTracker{})
	if err != nil {
		return nil, err
	}

	return &RaftDB{
		raft: raft,
		db:   backingDB,
	}, nil
}

// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
func (r *RaftDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	msg := NewPostAlertsRaftLogMessage(alerts...)
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	r.raft.Apply(bytes, 0)
	return nil
}

// GetAlerts gets all the alerts currently in the database.
func (r *RaftDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return r.db.GetAlerts(ctx)
}
