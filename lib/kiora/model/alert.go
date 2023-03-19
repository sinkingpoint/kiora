package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// DEFAULT_TIMEOUT_INTERVAL is the length of time after first seeing an alert that we time out the alert
// if we haven't seen any other information about it
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
	case AlertStatusFiring, AlertStatusAcked, AlertStatusResolved, AlertStatusTimedOut:
		return true
	default:
		return false
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

	defaultTime := time.Time{}
	if a.StartTime == defaultTime {
		return errors.New("missing start time in alert")
	}

	if a.TimeOutDeadline != defaultTime && !a.TimeOutDeadline.After(a.StartTime) {
		return errors.New("timeout deadline is not after start time")
	}

	return nil
}

func (a *Alert) UnmarshalJSON(b []byte) error {
	rawAlert := struct {
		Labels          Labels            `json:"labels"`
		Annotations     map[string]string `json:"annotations"`
		Status          AlertStatus       `json:"status"`
		StartTime       time.Time         `json:"startTime"`
		EndTime         time.Time         `json:"endsAt"`
		TimeOutDeadline time.Time         `json:"timeOutDeadline,omitempty"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rawAlert); err != nil {
		return err
	}

	if rawAlert.StartTime.IsZero() {
		a.StartTime = time.Now()
	} else {
		a.StartTime = rawAlert.StartTime
	}

	if rawAlert.TimeOutDeadline.IsZero() {
		a.TimeOutDeadline = a.StartTime.Add(DEFAULT_TIMEOUT_INTERVAL)
	} else {
		a.TimeOutDeadline = rawAlert.TimeOutDeadline
	}

	if rawAlert.Status == AlertStatusResolved && rawAlert.EndTime.IsZero() {
		a.EndTime = time.Now()
	} else {
		a.EndTime = rawAlert.EndTime
	}

	a.Labels = rawAlert.Labels
	a.Annotations = rawAlert.Annotations
	a.Status = rawAlert.Status

	return a.validate()
}
