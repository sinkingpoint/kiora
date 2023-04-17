package config

import (
	"net/http"
	"text/template"

	"github.com/rs/zerolog"
)

type HTTPClientOpt interface{}

// NodeBus is the interface that nodes get passed to give them access to
// things they may need to function, like a logger, or an HTTP client.
type NodeBus interface {
	// HTTPClient returns an HTTP client that can be used to make HTTP requests.
	// The returned client should be configured with the given options.
	HTTPClient(opts ...HTTPClientOpt) *http.Client

	// Logger returns a logger that can be used to log messages.
	Logger(component string) zerolog.Logger

	// Template returns a template that can be used to render templates.
	Template(name string) *template.Template

	// RegisterTemplate registers a template with the bus that other nodes can reference.
	RegisterTemplate(name string, tmpl *template.Template) error
}
