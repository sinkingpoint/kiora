package raft

import (
	"io"

	"github.com/hashicorp/raft"
)

var _ raft.FSM = &alertTracker{}

// alertTracker is the raft interface that handles consensus for the state of
// alerts in the system.
type alertTracker struct{}

func (a *alertTracker) Apply(*raft.Log) any {
	return nil
}

func (a *alertTracker) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (a *alertTracker) Restore(snapshot io.ReadCloser) error {
	return nil
}
