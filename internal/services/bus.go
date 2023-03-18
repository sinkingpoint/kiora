package services

import (
	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

// A bus wraps up all the things that a service might need to function.
type Bus interface {
	DB() kioradb.DB
	Broadcaster() clustering.Broadcaster
	Logger(serviceName string) *zerolog.Logger
}
