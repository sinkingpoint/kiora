package apiv1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_PROTO = "application/vnd.google.protobuf"

func Register(router *mux.Router, bus services.Bus, clusterer clustering.Clusterer) {
	api := apiv1{
		bus:       bus,
		clusterer: clusterer,
	}

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	subRouter.Path("/alerts").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(api.postAlerts), "POST api/v1/alerts"))
	subRouter.Path("/alerts").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(api.getAlerts), "GET /api/v1/alerts"))
	subRouter.Path("/alerts/ack").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(api.acknowledgeAlert), "POST /api/v1/alerts/ack"))
	subRouter.Path("/cluster/status").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(api.getClusterStatus), "GET /api/v1/cluster/status"))
	subRouter.Path("/silences").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(api.postSilences), "POST /api/v1/silences"))
	subRouter.Path("/silences").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(api.getSilences), "GET /api/v1/silences"))
}

type apiv1 struct {
	bus       services.Bus
	clusterer clustering.Clusterer
}

// postAlerts handles the POST /alerts request, decoding a list of alerts
// from the body, and forwarding them to the db.
func (a *apiv1) postAlerts(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var alerts []model.Alert

	switch r.Header.Get("Content-Type") {
	case CONTENT_TYPE_JSON:
		decoder := json.NewDecoder(io.NopCloser(bytes.NewBuffer(body)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&alerts); err != nil {
			span.RecordError(err)
			http.Error(w, fmt.Sprintf("failed to decode body: %q", err.Error()), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("invalid content-type %q", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
		return
	}

	if err := a.bus.Broadcaster().BroadcastAlerts(r.Context(), alerts...); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to process alerts")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// getAlerts returns all the alerts in the system as a JSON array.
// TODO(cdouch): Take filters here rather than returning everything.
func (a *apiv1) getAlerts(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	alerts := a.bus.DB().QueryAlerts(r.Context(), query.MatchAll())

	bytes, err := json.Marshal(alerts)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal alerts")
		log.Err(err).Msg("failed to get alerts")
		http.Error(w, "failed to get alerts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes) //nolint:errcheck
}

// getClusterStatus returns a JSON array of all the nodes in the cluster.
func (a *apiv1) getClusterStatus(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	if a.clusterer == nil {
		http.Error(w, "no clusterer configured", http.StatusNotFound)
		return
	}

	clusterNodes := a.clusterer.Nodes()
	bytes, err := json.Marshal(clusterNodes)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal cluster nodes")
		log.Err(err).Msg("failed to marshal cluster nodes")
		http.Error(w, "failed to marshal cluster nodes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (a *apiv1) acknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	type ackRequest struct {
		model.AlertAcknowledgement
		AlertID string `json:"alertID"`
	}

	ack := ackRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&ack); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // nolint:errcheck
		return
	}

	// TODO - validate this data.
	if err := a.bus.Broadcaster().BroadcastAlertAcknowledgement(r.Context(), ack.AlertID, ack.AlertAcknowledgement); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to broadcast alert acknowledgment")) // nolint:errcheck
		log.Err(err).Msg("failed to broadcast alert acknowledgment")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *apiv1) postSilences(w http.ResponseWriter, r *http.Request) {
	silence := model.Silence{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&silence); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // nolint:errcheck
		return
	}

	if err := a.bus.Broadcaster().BroadcastSilences(r.Context(), silence); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to broadcast silence")) // nolint:errcheck
		log.Err(err).Msg("failed to broadcast silence")
		return
	}

	responseBytes, _ := json.Marshal(silence) // TODO(cdouch): Error checking.

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes) // nolint:errcheck
}

func (a *apiv1) getSilences(w http.ResponseWriter, r *http.Request) {
	silences := a.bus.DB().QuerySilences(r.Context(), query.MatchAll())

	responseBytes, _ := json.Marshal(silences) // TODO(cdouch): Error checking.

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes) // nolint:errcheck
}
