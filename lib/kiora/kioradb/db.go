package kioradb

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
)

// ModelReader defines an interface that can be used to get various models out of an underlying system.
type ModelReader interface {
	// GetAlerts gets all the alerts currently in the database.
	GetAlerts(ctx context.Context) ([]model.Alert, error)

	// GetSilence returns all the silences that would match a given label set
	GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error)

	QueryAlerts(ctx context.Context, query AlertQuery) []model.Alert
}

// ModelWriter defines an interface that takes various models and processes them, with no way to get them back out.
type ModelWriter interface {
	// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
	ProcessAlerts(ctx context.Context, alerts ...model.Alert) error
	// ProcessSilences takes silences and processes them.
	ProcessSilences(ctx context.Context, silences ...model.Silence) error
}

// DB defines an interface that is able to process alerts, silences etc and store them (for some definition of store).
type DB interface {
	ModelReader
	ModelWriter

	RegisterEndpoints(ctx context.Context, httpRouter *mux.Router, grpcServer *grpc.Server) error
}

type FallthroughDB struct {
	db DB
}

func NewFallthroughDB(db DB) *FallthroughDB {
	return &FallthroughDB{
		db: db,
	}
}

// GetAlerts gets all the alerts currently in the database.
func (f *FallthroughDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	return f.db.GetAlerts(ctx)
}

// GetExistingAlert returns the existing data for a given labelset, if it exists.
func (f *FallthroughDB) QueryAlerts(ctx context.Context, query AlertQuery) []model.Alert {
	return f.db.QueryAlerts(ctx, query)
}

// GetSilence returns all the silences that would match a given label set
func (f *FallthroughDB) GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error) {
	return f.db.GetSilences(ctx, labels)
}

// ProcessAlerts takes alerts and processes them, adding new ones and resolving old ones.
func (f *FallthroughDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	return f.db.ProcessAlerts(ctx, alerts...)
}

// ProcessSilences takes silences and processes them.
func (f *FallthroughDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	return f.db.ProcessSilences(ctx, silences...)
}

func (f *FallthroughDB) RegisterEndpoints(ctx context.Context, httpRouter *mux.Router, grpcServer *grpc.Server) error {
	return f.db.RegisterEndpoints(ctx, httpRouter, grpcServer)
}
