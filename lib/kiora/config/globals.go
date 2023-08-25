package config

import (
	"net/http"
	"text/template"

	"github.com/rs/zerolog"
)

type HTTPClientOpt interface{}

// Globals is the interface that nodes get passed to give them access to
// things they may need to function, like a logger, or an HTTP client.
type Globals struct {
	httpClient *http.Client
	logger     zerolog.Logger
	templates  *template.Template

	Tenanter Tenanter
}

// GlobalsOpt is the interface that can be passed to NewGlobals to configure
// the globals object.
type GlobalsOpt func(*Globals)

// WithHTTPClient sets the HTTP client that will be returned by HTTPClient.
func WithHTTPClient(client *http.Client) GlobalsOpt {
	return func(g *Globals) {
		g.httpClient = client
	}
}

// WithLogger sets the logger that will be returned by Logger.
func WithLogger(logger zerolog.Logger) GlobalsOpt {
	return func(g *Globals) {
		g.logger = logger
	}
}

// WithTemplates sets the templates that will be returned by Template.
func WithTemplates(templates *template.Template) GlobalsOpt {
	return func(g *Globals) {
		g.templates = templates
	}
}

// WithTenanter sets the tenanter that will be returned by Tenanter.
func WithTenanter(t Tenanter) GlobalsOpt {
	return func(g *Globals) {
		g.Tenanter = t
	}
}

// NewGlobals creates a new Globals object with the given options.
func NewGlobals(opts ...GlobalsOpt) *Globals {
	g := &Globals{
		httpClient: http.DefaultClient,
		logger:     zerolog.Nop(),
		templates:  template.New(""),
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// HTTPClient returns an HTTP client that can be used to make HTTP requests.
// The returned client should be configured with the given options.
func (g *Globals) HTTPClient(opts ...HTTPClientOpt) *http.Client {
	return g.httpClient
}

// Logger returns a logger that can be used to log messages.
func (g *Globals) Logger(component string) zerolog.Logger {
	return g.logger.With().Str("component", component).Logger()
}

// Template returns a template that can be used to render templates.
func (g *Globals) Template(name string) *template.Template {
	return g.templates.Lookup(name)
}

// RegisterTemplate registers a template with the bus that other nodes can reference.
func (g *Globals) RegisterTemplate(name string, tmpl *template.Template) error {
	_, err := g.templates.AddParseTree(name, tmpl.Tree)
	return err
}
