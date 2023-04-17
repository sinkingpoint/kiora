package config

import (
	"net/http"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

var _ = config.NodeBus(&KioraNodeBus{})

type KioraNodeBus struct {
	logger    zerolog.Logger
	templates *template.Template
}

func NewKioraNodeBus(config *configGraph, logger zerolog.Logger) *KioraNodeBus {
	return &KioraNodeBus{
		logger:    logger,
		templates: template.New("kiora"),
	}
}

func (k *KioraNodeBus) HTTPClient(opts ...config.HTTPClientOpt) *http.Client {
	return &http.Client{}
}

func (k *KioraNodeBus) Logger(component string) zerolog.Logger {
	return k.logger.With().Str("component", component).Logger()
}

func (k *KioraNodeBus) Template(name string) *template.Template {
	return k.templates.Lookup(name)
}

func (k *KioraNodeBus) RegisterTemplate(name string, tmpl *template.Template) error {
	_, err := k.templates.AddParseTree(name, tmpl.Tree)
	return err
}
