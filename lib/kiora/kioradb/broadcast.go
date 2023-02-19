package kioradb

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
)

type Broadcaster interface {
	RegisterEndpoints(ctx context.Context, router *mux.Router, grcpServer *grpc.Server) error
	BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error
}
