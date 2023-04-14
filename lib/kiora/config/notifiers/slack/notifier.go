package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func init() {
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

	apiURL *config.MaybeSecretFile
}

func New(name string, bus config.NodeBus, attrs map[string]string) (config.Node, error) {
	if attrs["api_url"] != "" && attrs["api_url_file"] != "" {
		return nil, errors.New("cannot specify both api_url and api_url_file")
	}

	if attrs["api_url"] == "" && attrs["api_url_file"] == "" {
		return nil, errors.New("must specify either api_url or api_url_file")
	}

	apiURL, err := config.NewMaybeSecretFile(attrs["api_url_file"], config.Secret(attrs["api_url"]))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load api url")
	}

	return &SlackNotifier{
		name:   config.NotifierName(name),
		bus:    bus,
		client: bus.HTTPClient(),

		apiURL: apiURL,
	}, nil
}

func (s *SlackNotifier) Name() config.NotifierName {
	return s.name
}

func (s *SlackNotifier) Type() string {
	return "slack"
}

func (s *SlackNotifier) Notify(ctx context.Context, alerts ...model.Alert) *config.NotificationError {
	payload := slackPayload{
		Text: fmt.Sprintf("[FIRING: %d] %s", len(alerts), alerts[0].Labels["alertname"]),
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
