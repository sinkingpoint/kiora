package raft

// This contains the main raft state machine. You'll see a bunch of panics here. That is on purpose. Basically any parsing error we get here,
// we panic to avoid generating an inconsistency in the states between instances.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/protobuf/proto"
)

var _ raft.FSM = &kioraFSM{}
var _ raft.FSMSnapshot = &kioraSnapshot{}

type kioraSnapshot struct {
	Alerts []model.Alert `json:"alerts"`
}

func (k *kioraSnapshot) Persist(sink raft.SnapshotSink) error {
	bytes, err := json.Marshal(k)
	if err != nil {
		return err
	}

	_, err = sink.Write(bytes)
	return err
}

func (k *kioraSnapshot) Release() {}

// kioraFSM is the raft interface that handles consensus for the state of alerts in the system.
type kioraFSM struct {
	db kioradb.DB
}

func (a *kioraFSM) Apply(l *raft.Log) any {
	ctx, span := tracing.Tracer().Start(context.Background(), "kioraFSM.Apply")
	defer span.End()

	log, err := decodeLogMessage(l.Data)
	if err != nil {
		panic(fmt.Sprintf("BUG: failed to unmarshal raft message (%q). Stopping to avoid an inconsistency. This should never happen, please report.", err))
	}

	switch msg := log.Log.(type) {
	case *kioraproto.KioraLogMessage_Alerts:
		a.processAlerts(ctx, msg.Alerts)
	default:
		panic(fmt.Sprintf("BUG: Got a type of message that we haven't handled (%q)", msg))
	}

	return nil
}

// processAlerts handles the Alerts raft message, decoding the alerts into the model
// and passing them into the db for further processing.
func (a *kioraFSM) processAlerts(ctx context.Context, protoAlerts *kioraproto.PostAlertsMessage) {
	alerts := []model.Alert{}

	for _, protoAlert := range protoAlerts.Alerts {
		alert := model.Alert{}

		if err := alert.DeserializeFromProto(protoAlert); err != nil {
			panic(fmt.Sprintf("BUG: failed to unmarshal a model.Alert from a proto alert: %q", err))
		}
		alerts = append(alerts, alert)
	}

	if err := a.db.StoreAlerts(ctx, alerts...); err != nil {
		panic(fmt.Sprintf("BUG: failed to process alerts: %q", err))
	}
}

func (a *kioraFSM) Snapshot() (raft.FSMSnapshot, error) {
	alerts := a.db.QueryAlerts(context.Background(), &kioradb.AllMatchQuery{})

	return &kioraSnapshot{
		Alerts: alerts,
	}, nil
}

func (a *kioraFSM) Restore(input io.ReadCloser) error {
	decoder := json.NewDecoder(input)
	decoder.DisallowUnknownFields()
	var snapshot kioraSnapshot
	if err := decoder.Decode(&snapshot); err != nil {
		return err
	}

	return a.db.StoreAlerts(context.Background(), snapshot.Alerts...)
}

// decodeLogMessage decodes the raw bytes into a kioraproto.RaftLog.
func decodeLogMessage(raw []byte) (*kioraproto.KioraLogMessage, error) {
	msg := kioraproto.KioraLogMessage{}

	if err := proto.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
