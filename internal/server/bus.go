package server

import (
	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

type KioraBus struct {
	db          kioradb.DB
	broadcaster clustering.Broadcaster
	logger      *zerolog.Logger
	config      config.Config
}

func NewKioraBus(db kioradb.DB, broadcaster clustering.Broadcaster, config config.Config) *KioraBus {
	return &KioraBus{
		db:          db,
		broadcaster: broadcaster,
		logger:      zerolog.DefaultContextLogger,
		config:      config,
	}
}

func (k *KioraBus) DB() kioradb.DB {
	return k.db
}

func (k *KioraBus) Broadcaster() clustering.Broadcaster {
	return k.broadcaster
}

func (k *KioraBus) Logger(serviceName string) *zerolog.Logger {
	logger := k.logger.With().Str("service_name", serviceName).Logger()
	return &logger
}

func (k *KioraBus) Config() config.Config {
	return k.config
}
