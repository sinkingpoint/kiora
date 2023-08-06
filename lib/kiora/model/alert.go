package model

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/internal/stubs"
)

// DEFAULT_TIMEOUT_INTERVAL is the length of time after first seeing an alert that we time out the alert
// if we haven't seen any other information about it.
const DEFAULT_TIMEOUT_INTERVAL = 12 * time.Hour

// AlertStatus is the current status of an alert in Kiora.
type AlertStatus string

const (
	// AlertStatusFiring marks alerts that are currently active.
	AlertStatusFiring AlertStatus = "firing"

	// AlertStatusAcked marks alerts that are firing, but have been acknowledged by a human.
	AlertStatusAcked AlertStatus = "acked"

	// AlertStatusResolved marks alerts that were firing but have now been resolved.
	AlertStatusResolved AlertStatus = "resolved"

	// AlertStatusTimedOut marks alerts that we never got a resolved notification for, but hit their expiry times.
	AlertStatusTimedOut AlertStatus = "timed out"

	// AlertStatusSilenced marks alerts that have been silenced by one or more silences.
	AlertStatusSilenced AlertStatus = "silenced"
)

func (s AlertStatus) isValid() bool {
	switch s {
	case AlertStatusFiring, AlertStatusAcked, AlertStatusResolved, AlertStatusTimedOut, AlertStatusSilenced:
		return true
	default:
		return false
	}
}

// Alert is the _operational state_ of the alert. As opposed to the protobuf structs
// that are the values being transmitted, this struct contains all the state that might
// be ascertained by Kiora through interactions with other models (e.g. silences).
type Alert struct {
	// ID is the unique ID of the alert.
	ID string `json:"id,omitempty"`

	// Labels defines the metadata on the alert that is used for deduplication purposes.
	Labels Labels `json:"labels"`

	// Annotations defines them metadata on the alert that _isn't_ used for deduplication. This can be links etc.
	Annotations map[string]string `json:"annotations"`

	// Status is the status of the alert in the system.
	Status AlertStatus `json:"status"`

	// Acknowledgement is the details if this alert has fired and been acknowledged.
	Acknowledgement *AlertAcknowledgement `json:"acknowledgement,omitempty"`

	// StartTime is when the alert first started firing.
	StartTime time.Time `json:"startsAt"`

	// EndTime is when the alert ended (either timed out or resolved).
	EndTime time.Time `json:"endsAt"`

	// TimeOutDeadline is when the alert should be marked as timed out, assuming no further messages come in.
	TimeOutDeadline time.Time `json:"timeOutDeadline,omitempty"`

	// LastNotifyTime is the time that a notification for this alert was last sent.
	LastNotifyTime time.Time `json:"-"`
}

func (a *Alert) validate() error {
	if a.Labels == nil {
		return errors.New("missing labels in alert")
	}

	if a.Annotations == nil {
		return errors.New("missing annotations in alert")
	}

	if !a.Status.isValid() {
		return fmt.Errorf("invalid alert status in alert: %q", a.Status)
	}

	if a.StartTime.IsZero() {
		return errors.New("missing start time in alert")
	}

	if !a.EndTime.IsZero() && a.EndTime.Before(a.StartTime) {
		return errors.New("end time is before start time")
	}

	if !a.TimeOutDeadline.IsZero() && a.TimeOutDeadline.Before(a.StartTime) {
		return errors.New("timeout deadline is not after start time")
	}

	return nil
}

func (a *Alert) UnmarshalJSON(b []byte) error {
	rawAlert := struct {
		ID              string                `json:"id"`
		Labels          Labels                `json:"labels"`
		Annotations     map[string]string     `json:"annotations"`
		Status          AlertStatus           `json:"status"`
		StartTime       time.Time             `json:"startsAt"`
		EndTime         time.Time             `json:"endsAt"`
		TimeOutDeadline time.Time             `json:"timeOutDeadline,omitempty"`
		Acknowledgement *AlertAcknowledgement `json:"acknowledgement"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rawAlert); err != nil {
		return errors.Wrap(err, "failed to decode alert")
	}

	a.Labels = rawAlert.Labels
	a.Annotations = rawAlert.Annotations
	a.Status = rawAlert.Status
	a.StartTime = rawAlert.StartTime
	a.EndTime = rawAlert.EndTime
	a.Acknowledgement = rawAlert.Acknowledgement
	a.TimeOutDeadline = rawAlert.TimeOutDeadline

	return a.Materialise()
}

// Materialise fills in any missing fields in the alert with sensible defaults.
func (a *Alert) Materialise() error {
	if a.StartTime.IsZero() {
		a.StartTime = stubs.Time.Now()
	}

	if a.Annotations == nil {
		a.Annotations = map[string]string{}
	}

	if a.Status == AlertStatusResolved && a.EndTime.IsZero() {
		a.EndTime = stubs.Time.Now()
	}

	if a.TimeOutDeadline.IsZero() {
		a.TimeOutDeadline = a.StartTime.Add(DEFAULT_TIMEOUT_INTERVAL)
	}

	// AlertIDs are a bit arbitrary, but having them as a hash of the labels affords a few nice advantages.
	// Namely, it means that any given alert has a consistent ID across all Kiora instances, and across time.
	a.ID = alertID(a.Labels)
	return a.validate()
}

// Acknowledge marks this alert as Acknowledged with the given metadata.
func (a *Alert) Acknowledge(ack *AlertAcknowledgement) error {
	if a.Status != AlertStatusFiring {
		return errors.New("cannot acknowledge a non-firing alert")
	}

	a.Status = AlertStatusAcked
	a.Acknowledgement = ack
	return nil
}

func (a *Alert) Field(name string) (any, error) {
	if val, ok := a.Labels[name]; ok {
		return val, nil
	}

	// Special case non-label fields to allow filtering, and sorting on them.
	switch name {
	case "__id__":
		return a.ID, nil
	case "__status__":
		return a.Status, nil
	case "__starts_at__":
		return a.StartTime, nil
	case "__ends_at__":
		return a.EndTime, nil
	case "__timeout_deadline__":
		return a.TimeOutDeadline, nil
	case "__last_notify_time__":
		return a.LastNotifyTime, nil
	}

	return "", fmt.Errorf("label %q doesn't exist", name)
}

// alertID is a helper function to generate an AlertID from a labelset.
func alertID(labels Labels) string {
	bytes := make([]byte, 8) // 8 bytes in a uint64.
	binary.LittleEndian.PutUint64(bytes, labels.Hash())
	return hex.EncodeToString(bytes)
}
