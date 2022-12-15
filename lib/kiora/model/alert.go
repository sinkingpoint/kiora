package model

import (
	"fmt"
	"time"

	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
)

// DEFAULT_TIMEOUT_INTERVAL is the length of time after first seeing an alert that we time out the alert
// if we haven't seen any other information about it
const DEFAULT_TIMEOUT_INTERVAL = 12 * time.Hour

// AlertStatus is the current status of an alert in Kiora.
type AlertStatus string

const (
	// AlertStatusFiring marks alerts that are currently active.
	AlertStatusFiring AlertStatus = "firing"

	// AlertStatusProcessing marks alerts that have been accepted, but aren't active for whatever reason.
	AlertStatusProcessing AlertStatus = "processing"

	// AlertStatusAcked marks alerts that are firing, but have been acknowledged by a human.
	AlertStatusAcked AlertStatus = "acked"

	// AlertStatusResolved marks alerts that were firing but have now been resolved.
	AlertStatusResolved AlertStatus = "resolved"

	// AlertStatusTimedOut marks alerts that we never got a resolved notification for, but hit their expiry times.
	AlertStatusTimedOut AlertStatus = "timed out"
)

// deserializeStatusFromProto takes the proto AlertStatus and turns it into a model.AlertStatus
func deserializeStatusFromProto(status kioraproto.AlertStatus) AlertStatus {
	switch status {
	case kioraproto.AlertStatus_firing:
		return AlertStatusFiring
	case kioraproto.AlertStatus_resolved:
		return AlertStatusResolved
	default:
		panic(fmt.Sprintf("BUG: unhandled alert status received from proto: %q", status.String()))
	}
}

// Alert is the _operational state_ of the alert. As opposed to the protobuf structs
// that are the values being transmitted, this struct contains all the state that might
// be ascertained by Kiora through interactions with other models (e.g. silences).
type Alert struct {
	// Labels defines the metadata on the alert that is used for deduplication purposes.
	Labels Labels `json:"labels"`

	// Annotations defines them metadata on the alert that _isn't_ used for deduplication. This can be links etc.
	Annotations map[string]string `json:"annotations"`

	// Status is the status of the alert in the system.
	Status AlertStatus `json:"status"`

	// StartTime is when the alert first started firing.
	StartTime time.Time `json:"startTime"`

	// TimeOutDeadline is when the alert should be marked as timed out, assuming no further messages come in.
	TimeOutDeadline time.Time `json:"timeOutDeadline,omitempty"`
}

// DeserializeFromProto creates a model.Alert from a proto alert
func (a *Alert) DeserializeFromProto(proto *kioraproto.Alert) error {
	a.Labels = proto.Labels
	a.Annotations = proto.Annotations
	a.Status = deserializeStatusFromProto(proto.Status)

	if proto.StartTime != nil {
		a.StartTime = time.UnixMilli(proto.StartTime.AsTime().UnixMilli()).UTC()
	}

	if proto.EndTime != nil && proto.EndTime.Nanos > 0 {
		a.TimeOutDeadline = time.UnixMilli(proto.EndTime.AsTime().UnixMilli())
	} else {
		a.TimeOutDeadline = a.StartTime.Add(DEFAULT_TIMEOUT_INTERVAL)
	}

	return nil
}
