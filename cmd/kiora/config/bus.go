package config

import (
	"net/http"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

type KioraNodeBus struct {
}

func NewKioraNodeBus(config *configGraph) *KioraNodeBus {
	return &KioraNodeBus{}
}

func (k *KioraNodeBus) HTTPClient(opts ...config.HTTPClientOpt) *http.Client {
	return &http.Client{}
}
