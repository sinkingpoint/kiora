package kiora

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type KioraInstance struct {
	HTTPClient *http.Client
	APIVersion string
	URL        string
}

func NewKioraInstance(url string, apiVersion string) *KioraInstance {
	return &KioraInstance{
		HTTPClient: http.DefaultClient,
		APIVersion: apiVersion,
		URL:        url,
	}
}

func (k *KioraInstance) getRequest(method string, uri string) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/%s/%s", k.URL, k.APIVersion, uri)
	return http.NewRequest(method, url, nil)
}

func (k *KioraInstance) GetAlerts() ([]model.Alert, error) {
	req, err := k.getRequest(http.MethodGet, "alerts")
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
