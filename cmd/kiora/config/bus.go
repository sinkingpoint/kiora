package config

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

var _ = config.NodeBus(&KioraNodeBus{})

type KioraNodeBus struct {
	logger zerolog.Logger
}

func NewKioraNodeBus(config *configGraph, logger zerolog.Logger) *KioraNodeBus {
	return &KioraNodeBus{
		logger: logger,
	}
}

func (k *KioraNodeBus) HTTPClient(opts ...config.HTTPClientOpt) *http.Client {
	return &http.Client{}
}

func (k *KioraNodeBus) Logger(component string) zerolog.Logger {
	return k.logger.With().Str("component", component).Logger()
}
