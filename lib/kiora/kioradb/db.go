package kioradb

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// DB defines an interface that is able to process alerts, silences etc and store them (for some definition of store).
type DB interface {
	// StoreAlerts stores the given alerts in the database, overriding any existing alerts with the same labels.
	StoreAlerts(ctx context.Context, alerts ...model.Alert) error

	// QueryAlerts queries the database for alerts matching the given query.
	QueryAlerts(ctx context.Context, query query.AlertQuery) []model.Alert
}
