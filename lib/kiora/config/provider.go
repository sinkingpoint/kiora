package config

import (
	"context"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type NotifierName string

// DefaultGroupWait is the default amount of time that alerts sit around waiting for more alerts to be added to the group.
// This is pretty arbitrary, but increasing it will increase the amount of time that alerts are delayed, while sending fewer
// notifications. Decreasing it will decrease the amount of time that alerts are delayed, but will send more notifications.
const DEFAULT_GROUP_WAIT = 10 * time.Second

// NotificationError represents an error that occurred while sending a notification.
type NotificationError struct {
	Err       error
	Retryable bool
}

func NewNotificationError(err error, retryable bool) *NotificationError {
	return &NotificationError{
		Err:       err,
		Retryable: retryable,
	}
}

func (n *NotificationError) Error() string {
	return n.Err.Error()
}

// Notifier represents something that can send a notification about an alert.
type Notifier interface {
	// Name returns the name of the notifier.
	Name() NotifierName

	// Notify sends a notification about the given alerts.
	Notify(ctx context.Context, alerts ...model.Alert) *NotificationError
}

// Config represents a configuration that can return a list of notifiers for a given alert.
type Config interface {
	// Returns the notifiers that should be invoked for the given alert. If the response is nil,
	// then the notifier should do nothing, as opposed to an empty array that represents that the alert
	// should be processed as if it should be considered to be properly notified.
	GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []NotifierSettings

	// ValidateData returns an error that can be displayed to the user if the
	// data is invalid according to whatever rules the config has.
	ValidateData(ctx context.Context, data Fielder) error
}

// NotifierSettings represents a Notifier with additional settings. Such as grouping, and rate limiting settings.
type NotifierSettings struct {
	Notifier

	// GroupLabels is a list of label names that should be used to group alerts together.
	GroupLabels []string

	// GroupWait is the amount of time to wait before sending a notification for a group of alerts, to give time for more alerts to be added to the group.
	GroupWait time.Duration
}

func DefaultNotifierSettings() NotifierSettings {
	return NotifierSettings{
		GroupLabels: []string{"alertname"},
		GroupWait:   DEFAULT_GROUP_WAIT,
	}
}

// NewNotifier creates a new NotifierSettings with the given Notifier, and default settings.
func NewNotifier(n Notifier) NotifierSettings {
	return DefaultNotifierSettings().WithNotifier(n)
}

func (n NotifierSettings) WithGroupLabels(labels ...string) NotifierSettings {
	n.GroupLabels = labels
	return n
}

func (n NotifierSettings) WithGroupWait(wait time.Duration) NotifierSettings {
	n.GroupWait = wait
	return n
}

func (n NotifierSettings) WithNotifier(notifier Notifier) NotifierSettings {
	n.Notifier = notifier
	return n
}
