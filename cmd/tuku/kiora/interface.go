package kiora

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type KioraInstance struct {
	HTTPClient *http.Client
	APIVersion string
	URL        string
}

func NewKioraInstance(url, apiVersion string) *KioraInstance {
	return &KioraInstance{
		HTTPClient: http.DefaultClient,
		APIVersion: apiVersion,
		URL:        url,
	}
}

func (k *KioraInstance) getRequest(method, uri string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/%s/%s", k.URL, k.APIVersion, uri)
	return http.NewRequest(method, url, body)
}

func (k *KioraInstance) GetAlerts() ([]model.Alert, error) {
	req, err := k.getRequest(http.MethodGet, "alerts", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	resp, err := k.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	alerts := []model.Alert{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&alerts); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return alerts, nil
}

func (k *KioraInstance) PostAlerts(alerts []model.Alert) error {
	body, err := json.Marshal(alerts)
	if err != nil {
		return errors.Wrap(err, "failed to marshal alerts")
	}

	req, err := k.getRequest(http.MethodPost, "alerts", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := k.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d (%q)", resp.StatusCode, string(body))
	}

	return nil
}

func (k *KioraInstance) PostSilence(silence model.Silence) (*model.Silence, error) {
	body, err := json.Marshal(silence)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal silence")
	}

	req, err := k.getRequest(http.MethodPost, "silences", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := k.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d (%q)", resp.StatusCode, string(body))
	}

	silence = model.Silence{}
	if err := json.Unmarshal(body, &silence); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	return &silence, nil
}
