package config

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Notifier represents something that can send a notification about an alert.
type Notifier interface {
	Notify(ctx context.Context, alerts ...model.Alert) error
}

// Config represents a configuration that can return a list of notifiers for a given alert.
type Config interface {
	// Returns the notifiers that should be invoked for the given alert. If the response is nil,
	// then the notifier should do nothing, as opposed to an empty array that represents that the alert
	// should be processed as if it should be considered to be properly notified.
	GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []Notifier

	// ValidateData returns an error that can be displayed to the user if the
	// data is invalid according to whatever rules the config has.
	ValidateData(ctx context.Context, data Fielder) error
}
