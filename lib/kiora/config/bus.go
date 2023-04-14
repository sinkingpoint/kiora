package config

import "net/http"

type HTTPClientOpt interface{}

type NodeBus interface {
	HTTPClient(opts ...HTTPClientOpt) *http.Client
}
