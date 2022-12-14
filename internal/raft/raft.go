package raft

// This contains the main raft state machine. You'll see a bunch of panics here. That is on purpose. Basically any parsing error we get here,
// we panic to avoid generating an inconsistency in the states between instances.

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"capnproto.org/go/capnp/v3"
	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ raft.FSM = &alertTracker{}

// alertTracker is the raft interface that handles consensus for the state of
// alerts in the system.
type alertTracker struct {
	db kioradb.DB
}

func NewAlertTracker(db kioradb.DB) *alertTracker {
	return &alertTracker{
		db: db,
	}
}

func (a *alertTracker) Apply(l *raft.Log) any {
	msg, err := decodeLogMessage(l.Data)
	if err != nil {
		panic(fmt.Sprintf("BUG: failed to unmarshal raft message (%q). Stopping to avoid an inconsistency. This should never happen, please report.", err))
	}

	log := msg.Log()

	switch log.Which() {
	case kioraproto.RaftLog_log_Which_alerts:
		alertsWrapper, err := log.Alerts()
		if err != nil {
			panic(fmt.Sprintf("BUG: failed to unmarshal alerts (%q), in an alerts raft log. Stopping to avoid an inconsistency. This should never happen, please report.", err))
		}

		alerts, err := alertsWrapper.Alerts()
		if err != nil {
			panic(fmt.Sprintf("BUG: failed to unmarshal alerts (%q), in an alerts message. Stopping to avoid an inconsistency. This should never happen, please report.", err))
		}

		a.processAlerts(alerts)
	default:
		panic(fmt.Sprintf("BUG: Got a type of message that we haven't handled (%q)", log.Which()))
	}

	return nil
}

// processAlerts handles the Alerts raft message, decoding the alerts into the model
// and passing them into the db for further processing.
func (a *alertTracker) processAlerts(protoAlerts kioraproto.Alert_List) {
	alerts := []model.Alert{}

	for i := 0; i < protoAlerts.Len(); i++ {
		protoAlert := protoAlerts.At(i)
		var alert model.Alert
		if err := alert.DeserializeFromProto(&protoAlert); err != nil {
			panic(fmt.Sprintf("BUG: failed to unmarshal a model.Alert from a proto alert: %q", err))
		}

		alerts = append(alerts, alert)
	}

	if err := a.db.ProcessAlerts(context.Background(), alerts...); err != nil {
		panic(fmt.Sprintf("BUG: failed to process alerts: %q", err))
	}
}

func (a *alertTracker) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (a *alertTracker) Restore(snapshot io.ReadCloser) error {
	return nil
}

// decodeLogMessage decodes the raw bytes into a kioraproto.RaftLog
func decodeLogMessage(raw []byte) (kioraproto.RaftLog, error) {
	decoder := capnp.NewDecoder(io.NopCloser(bytes.NewBuffer(raw)))
	msg, err := decoder.Decode()
	if err != nil {
		return kioraproto.RaftLog{}, err
	}

	return kioraproto.ReadRootRaftLog(msg)
}
