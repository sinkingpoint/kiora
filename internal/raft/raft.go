package raft

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const LOGDB_FILE_NAME = "log.dat"
const STABLEDB_PATH = "stable.dat"
const SNAPSHOTDB_FILE_NAME = "snapshot.dat"

type raftConfig struct {
	LocalID           string
	LocalAddress      string
	DataDir           string
	SnapshotRetention int
	Bootstrap         bool
}

func NewRaftConfig(localID string) raftConfig {
	return raftConfig{
		LocalID:           localID,
		LocalAddress:      "localhost:4279",
		DataDir:           "./kiora/data",
		SnapshotRetention: 3,
		Bootstrap:         true,
	}
}

func NewRaft(ctx context.Context, config raftConfig, stateMachine *alertTracker) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(config.LocalID)
	baseDir := filepath.Join(config.DataDir, config.LocalID)
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return nil, nil, err
	}

	logDBPath := filepath.Join(baseDir, LOGDB_FILE_NAME)
	logDB, err := boltdb.NewBoltStore(logDBPath)
	if err != nil {
		return nil, nil, err
	}

	stableDBStorePath := filepath.Join(baseDir, SNAPSHOTDB_FILE_NAME)
	stableDB, err := boltdb.NewBoltStore(stableDBStorePath)
	if err != nil {
		return nil, nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, config.SnapshotRetention, os.Stderr)
	if err != nil {
		return nil, nil, err
	}

	// TODO(cdouch): Allow securing the transport with the config.
	tm := transport.New(raft.ServerAddress(config.LocalAddress), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})

	r, err := raft.NewRaft(c, stateMachine, logDB, stableDB, snapshotStore, tm.Transport())
	if err != nil {
		return nil, nil, err
	}

	if config.Bootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(config.LocalID),
					Address:  raft.ServerAddress(config.LocalAddress),
				},
			},
		}

		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return r, tm, nil
}
