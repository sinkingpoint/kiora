package apiv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/internal/server/api"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_PROTO = "application/vnd.google.protobuf"

func Register(router *mux.Router, api api.API, logger zerolog.Logger) {
	apiv1 := apiv1{
		api:    api,
		logger: logger.With().Str("component", "apiv1").Logger(),
	}

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	subRouter.Path("/alerts").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.postAlerts), "POST api/v1/alerts"))
	subRouter.Path("/alerts").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.getAlerts), "GET /api/v1/alerts"))
	subRouter.Path("/alerts/stats").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.queryAlertStats), "GET /api/v1/alerts/stats"))
	subRouter.Path("/alerts/ack").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.acknowledgeAlert), "POST /api/v1/alerts/ack"))
	subRouter.Path("/cluster/status").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.getClusterStatus), "GET /api/v1/cluster/status"))
	subRouter.Path("/silences").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.postSilence), "POST /api/v1/silences"))
	subRouter.Path("/silences").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.getSilences), "GET /api/v1/silences"))
}

type apiv1 struct {
	api    api.API
	logger zerolog.Logger
}

func New(api api.API, logger zerolog.Logger) *apiv1 {
	return &apiv1{
		api:    api,
		logger: logger.With().Str("component", "apiv1").Logger(),
	}
}

func decodeFromContentType(r *http.Request, v interface{}) error {
	_, span := otel.Tracer("").Start(r.Context(), "apiv1.decodeFromContentType")
	defer span.End()

	switch r.Header.Get("Content-Type") {
	case CONTENT_TYPE_JSON:
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		return decoder.Decode(v)
	case CONTENT_TYPE_PROTO:
		return nil
	default:
		return fmt.Errorf("unsupported content type: %q", r.Header.Get("Content-Type"))
	}
}

// postAlerts handles the POST /alerts request, decoding a list of alerts
// from the body, and forwarding them to the db.
func (a *apiv1) postAlerts(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	var alerts []model.Alert

	switch r.Header.Get("Content-Type") {
	case CONTENT_TYPE_JSON:
		if err := decodeFromContentType(r, &alerts); err != nil {
			span.RecordError(err)
			http.Error(w, fmt.Sprintf("failed to decode body: %q", err.Error()), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("invalid content-type %q", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
		return
	}

	if err := a.api.PostAlerts(r.Context(), alerts); err != nil {
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

	var queries []query.AlertFilter
	if id := r.URL.Query().Get("id"); id != "" {
		queries = append(queries, query.ID(id))
	}

	opts := []query.QueryOp{}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			span.RecordError(err)
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		opts = append(opts, query.Limit(limit))
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			span.RecordError(err)
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		opts = append(opts, query.Offset(offset))
	}

	if sort := r.URL.Query()["sort"]; len(sort) > 0 {
		order := query.OrderAsc
		if r.URL.Query().Get("order") == string(query.OrderDesc) {
			order = query.OrderDesc
		}

		opts = append(opts, query.OrderBy(sort, order))
	}

	var q query.AlertFilter
	if len(queries) == 0 {
		q = query.MatchAll()
	} else {
		q = query.AllAlerts(queries...)
	}

	alerts, err := a.api.GetAlerts(r.Context(), query.NewAlertQuery(q, opts...))
	if err != nil {
		span.RecordError(err)
		http.Error(w, "failed to get alerts", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(alerts)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "failed to get alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) //nolint:errcheck
}

func (a *apiv1) queryAlertStats(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	args := make(map[string]string)
	for k, v := range r.URL.Query() {
		args[k] = v[0]
	}

	q, err := query.UnmarshalAlertStatsQuery(args)
	if err != nil {
		a.logger.Warn().Msg(err.Error())
		span.RecordError(err)
		http.Error(w, "failed to unmarshal query", http.StatusBadRequest)
		return
	}

	stats, err := a.api.QueryAlertStats(r.Context(), q)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "failed to query alert stats", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(stats)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "failed to marshal alert stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) //nolint:errcheck
}

// getClusterStatus returns a JSON array of all the nodes in the cluster.
func (a *apiv1) getClusterStatus(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	clusterNodes, err := a.api.GetClusterStatus(r.Context())
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if clusterNodes == nil {
		http.Error(w, "no clusterer configured", http.StatusNotFound)
	}

	bytes, err := json.Marshal(clusterNodes)

	if err != nil {
		span.SetStatus(codes.Error, "failed to marshal cluster nodes")
		a.logger.Err(err).Msg("failed to marshal cluster nodes")
		http.Error(w, "failed to marshal cluster nodes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
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

	if ack.AlertID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("alertID is required")) // nolint:errcheck
		return
	}

	if err := a.api.AckAlert(r.Context(), ack.AlertID, ack.AlertAcknowledgement); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to handle alert acknowledgement: %s", err.Error()))) // nolint:errcheck
		a.logger.Err(err).Msg("failed to broadcast alert acknowledgment")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *apiv1) postSilence(w http.ResponseWriter, r *http.Request) {
	silence := model.Silence{}
	if err := decodeFromContentType(r, &silence); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // nolint:errcheck
		return
	}

	if err := a.api.PostSilence(r.Context(), silence); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to handle alert acknowledgement")) // nolint:errcheck
		a.logger.Err(err).Msg("failed to broadcast alert acknowledgment")
		return
	}

	responseBytes, err := json.Marshal(silence)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal silence")) // nolint:errcheck
		a.logger.Err(err).Msg("failed to marshal silence")
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes) // nolint:errcheck
}

func (a *apiv1) getSilences(w http.ResponseWriter, r *http.Request) {
	silences, err := a.api.GetSilences(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get silences")) // nolint:errcheck
		a.logger.Err(err).Msg("failed to get silences")
		return
	}

	responseBytes, err := json.Marshal(silences)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal silences")) // nolint:errcheck
		a.logger.Err(err).Msg("failed to marshal silences")
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes) // nolint:errcheck
}
