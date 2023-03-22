package promcompat

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/server/api"
	kmodel "github.com/sinkingpoint/kiora/lib/kiora/model"
)

func Register(router *mux.Router, api api.API) {
	promCompat := NewPromCompat(api)

	subRouter := router.PathPrefix("/api/prom-compat").Subrouter()

	subRouter.Path("/api/v2/alerts").Methods(http.MethodPost).HandlerFunc(promCompat.PostAlerts)
}

// promCompat provides an API that is able to ingest alerts from Prometheus.
type promCompat struct {
	api api.API
}

func NewPromCompat(api api.API) *promCompat {
	return &promCompat{
		api: api,
	}
}

// PostAlerts handles the POST /api/v2/alerts request, decoding a list of Prometheus alerts,
// converting them to Kiora alerts, and forwarding them to the db.
func (p *promCompat) PostAlerts(w http.ResponseWriter, r *http.Request) {
	promAlerts := []model.Alert{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&promAlerts); err != nil {
		http.Error(w, "failed to decode body", http.StatusBadRequest)
		return
	}

	alerts := make([]kmodel.Alert, len(promAlerts))
	for i, promAlert := range promAlerts {
		alert, err := marshalPromAlertToKioraAlert(promAlert)
		if err != nil {
			http.Error(w, "failed to unmarshal prometheus alert", http.StatusBadRequest)
			return
		}
		alerts[i] = alert
	}

	if err := p.api.PostAlerts(r.Context(), alerts); err != nil {
		http.Error(w, "failed to post alerts", http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to post alerts")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// marshalPromAlertToKioraAlert converts a Prometheus alert to a Kiora alert.
func marshalPromAlertToKioraAlert(p model.Alert) (kmodel.Alert, error) {
	labels := kmodel.Labels{}
	for k, v := range p.Labels {
		labels[string(k)] = string(v)
	}

	annotations := map[string]string{}
	for k, v := range p.Annotations {
		annotations[string(k)] = string(v)
	}

	alert := kmodel.Alert{
		Status:      kmodel.AlertStatus(p.Status()),
		Labels:      labels,
		Annotations: annotations,
		StartTime:   p.StartsAt,
		EndTime:     p.EndsAt,
	}

	return alert, alert.Materialise()
}
