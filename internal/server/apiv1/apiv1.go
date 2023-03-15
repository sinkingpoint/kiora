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
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_PROTO = "application/vnd.google.protobuf"

func Register(router *mux.Router, db kioradb.DB, broadcaster clustering.Broadcaster) {
	api := apiv1{
		db:          db,
		broadcaster: broadcaster,
	}
	router.Path("/api/v1/alerts").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(api.postAlerts), "POST api/v1/alerts"))
	router.Path("/api/v1/alerts").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(api.getAlerts), "GET /api/v1/alerts"))
}

type apiv1 struct {
	db          kioradb.DB
	broadcaster clustering.Broadcaster
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

	// For all incoming alerts, mark them as processing so later stages know that these are new alerts, not replayed old ones.
	for i := range alerts {
		if alerts[i].Status == model.AlertStatusFiring {
			alerts[i].Status = model.AlertStatusProcessing
		}
	}

	if err := a.broadcaster.BroadcastAlerts(r.Context(), alerts...); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to process alerts")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *apiv1) getAlerts(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	alerts := a.db.QueryAlerts(r.Context(), &kioradb.AllMatchQuery{})

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
