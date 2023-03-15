package serf

import "github.com/sinkingpoint/kiora/internal/clustering"

var _ = clustering.Broadcaster(&SerfBroadcaster{})

type Config struct{}

type SerfBroadcaster struct {
	conf *Config
}


func (s *SerfBroadcaster) RegisterEndpoints(ctx context.Context, router *mux.Router) error {
	
}
BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error