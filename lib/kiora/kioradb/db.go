package kioradb

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// ModelReader defines an interface that can be used to get various models out of an underlying system.
type ModelReader interface {
	// GetAlerts gets all the alerts currently in the database.
	GetAlerts(ctx context.Context) ([]model.Alert, error)

	// GetExistingAlert returns the existing data for a given labelset, if it exists.
	GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error)
}

type AlertProcessor interface {
	// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
	ProcessAlerts(ctx context.Context, alerts ...model.Alert) error
}

type SilenceProcessor interface {
	// ProcessSilences takes silences and processes them.
	ProcessSilences(ctx context.Context, silences ...model.Silence) error
}

// ModelWriter defines an interface that takes various models and processes them, with no way to get them back out.
type ModelWriter interface {
	AlertProcessor
	SilenceProcessor
}

// DB defines an interface that is able to process alerts, silences etc and store them (for some definition of store).
type DB interface {
	ModelReader
	ModelWriter
}
