package clustering

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/grpc"
)

type Server interface {
	Name() string
	Address() string
}

type Clusterer interface {
	GetMembers(ctx context.Context) ([]Server, error)
}

type Broadcaster interface {
	RegisterEndpoints(ctx context.Context, router *mux.Router, grcpServer *grpc.Server) error
	BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error
}
