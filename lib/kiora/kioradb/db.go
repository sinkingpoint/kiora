package kioradb

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// DB defines an interface that is able to process alerts, silences etc and store them (for some definition of store).
type DB interface {
	// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
	ProcessAlerts(ctx context.Context, alerts ...model.Alert) error
}
