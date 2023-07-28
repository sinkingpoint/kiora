package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var DefaultSlackTemplate *template.Template

func init() {
	tmpl, err := template.New("slack").Parse(`[FIRING: {{ len . }}] {{ (index . 0).Labels.alertname }}`)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse default slack template")
	}

	DefaultSlackTemplate = tmpl

	config.RegisterNode("slack", New)
}

type slackPayload struct {
	Text string `json:"text"`
}

// SlackNotifier is a notifier that sends alerts to a slack channel.
type SlackNotifier struct {
	name   config.NotifierName
	bus    config.NodeBus
	client *http.Client

	apiURL *unmarshal.MaybeSecretFile
}

func New(name string, bus config.NodeBus, attrs map[string]string) (config.Node, error) {
	rawNode := struct {
		ApiURL       *unmarshal.MaybeSecretFile `config:"api_url" required:"true"`
		TemplateFile *unmarshal.MaybeFile       `config:"template_file"`
	}{}

	if err := unmarshal.UnmarshalConfig(attrs, rawNode, unmarshal.UnmarshalOpts{
		DisallowUnknownFields: true,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	if err := bus.RegisterTemplate("slack", DefaultSlackTemplate); err != nil {
		return nil, err
	}

	return &SlackNotifier{
		name:   config.NotifierName(name),
		bus:    bus,
		client: bus.HTTPClient(),

		apiURL: rawNode.ApiURL,
	}, nil
}

func (s *SlackNotifier) Name() config.NotifierName {
	return s.name
}

func (s *SlackNotifier) Type() string {
	return "slack"
}

func (s *SlackNotifier) Notify(ctx context.Context, alerts ...model.Alert) *config.NotificationError {
	tmpl := s.bus.Template("slack")
	writer := strings.Builder{}
	if err := tmpl.Execute(&writer, alerts); err != nil {
		return &config.NotificationError{
			Err:       err,
			Retryable: false,
		}
	}

	payload := slackPayload{
		Text: writer.String(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return &config.NotificationError{
			Err:       err,
			Retryable: false,
		}
	}

	request, err := http.NewRequest(http.MethodPost, string(s.apiURL.Value()), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &config.NotificationError{
			Err:       err,
			Retryable: false,
		}
	}

	resp, err := s.client.Do(request)
	if err != nil {
		return &config.NotificationError{
			Err:       err,
			Retryable: true,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &config.NotificationError{
			Err:       errors.Errorf("unexpected status code: %d", resp.StatusCode),
			Retryable: true,
		}
	}

	return nil
}
