package apiv1

import (
	"encoding/json"
	"fmt"
	"net/http"

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

const (
	CONTENT_TYPE_JSON  = "application/json"
	CONTENT_TYPE_PROTO = "application/vnd.google.protobuf"
)

func Register(router *mux.Router, api api.API, logger zerolog.Logger) {
	baseAPI := &apiv1{
		api:    api,
		logger: logger.With().Str("component", "apiv1").Logger(),
	}

	apiv1 := ServerInterfaceWrapper{
		Handler: baseAPI,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	subRouter.Path("/alerts").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.PostAlerts), "POST api/v1/alerts"))
	subRouter.Path("/alerts").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.GetAlerts), "GET /api/v1/alerts"))
	subRouter.Path("/alerts/stats").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.GetAlertsStats), "GET /api/v1/alerts/stats"))
	subRouter.Path("/alerts/ack").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.PostAlertsAck), "POST /api/v1/alerts/ack"))
	subRouter.Path("/silences").Methods(http.MethodPost).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.PostSilences), "POST /api/v1/silences"))
	subRouter.Path("/silences").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(apiv1.GetSilences), "GET /api/v1/silences"))

	// This is technically not in the spec.
	subRouter.Path("/cluster/status").Methods(http.MethodGet).Handler(otelhttp.NewHandler(http.HandlerFunc(baseAPI.getClusterStatus), "GET /api/v1/cluster/status"))
}

var _ = ServerInterface(&apiv1{})

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
func (a *apiv1) PostAlerts(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	var body PostAlertsJSONBody

	switch r.Header.Get("Content-Type") {
	case CONTENT_TYPE_JSON:
		if err := decodeFromContentType(r, &body); err != nil {
			span.RecordError(err)
			http.Error(w, fmt.Sprintf("failed to decode body: %q", err.Error()), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("invalid content-type %q", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
		return
	}

	alerts := make([]model.Alert, len(body))
	for i, alert := range body {
		alerts[i] = model.Alert{
			Labels:      alert.Labels,
			Annotations: alert.Annotations,
			StartTime:   alert.StartsAt,
			Status:      model.AlertStatus(alert.Status),
		}

		if alert.EndsAt != nil {
			alerts[i].EndTime = *alert.EndsAt
		}

		if err := alerts[i].Materialise(); err != nil {
			span.RecordError(err)
			http.Error(w, fmt.Sprintf("failed to materialise alert: %q", err.Error()), http.StatusBadRequest)
			return
		}
	}

	if err := a.api.PostAlerts(r.Context(), alerts); err != nil {
		a.logger.Debug().Err(err).Msg("failed to post alerts")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func constructQueryOpts(limit *int, offset *int, sort *[]string, order string) ([]query.QueryOption, error) {
	opts := []query.QueryOption{}

	if limit != nil {
		opts = append(opts, query.Limit(*limit))
	}

	if offset != nil {
		opts = append(opts, query.Offset(*offset))
	}

	if sort != nil && len(*sort) > 0 {
		direction := query.OrderAsc
		if order == string(query.OrderDesc) {
			direction = query.OrderDesc
		}

		opts = append(opts, query.OrderBy(*sort, direction))
	} else if order != "" {
		return nil, fmt.Errorf("order specified without sort")
	}

	return opts, nil
}

// getAlerts returns all the alerts in the system as a JSON array.
func (a *apiv1) GetAlerts(w http.ResponseWriter, r *http.Request, params GetAlertsParams) {
	span := trace.SpanFromContext(r.Context())

	var queries []query.AlertFilter

	if params.Matchers != nil {
		for _, matcherString := range *params.Matchers {
			matcher := model.Matcher{}
			if err := matcher.UnmarshalText(matcherString); err != nil {
				a.logger.Debug().Err(err).Msgf("failed to unmarshal matcher %q", matcherString)
				span.RecordError(err)
				http.Error(w, fmt.Sprintf("failed to unmarshal matcher: %q", matcherString), http.StatusBadRequest)
				return
			}

			if matcher.Label == "__id__" && !matcher.IsRegex && !matcher.IsNegative {
				// For __id__=value matchers, we can use the ID filter, which is more efficient.
				queries = append(queries, query.ID(matcher.Value))
			} else {
				queries = append(queries, query.Matcher(matcher))
			}
		}
	}

	order := ""
	if params.Order != nil {
		order = string(*params.Order)
	}

	opts, err := constructQueryOpts(params.Limit, params.Offset, params.Sort, order)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to construct query options")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var q query.AlertFilter
	if len(queries) == 0 {
		q = query.MatchAll()
	} else {
		q = query.AllAlerts(queries...)
	}

	alerts, err := a.api.GetAlerts(r.Context(), query.NewAlertQuery(q, opts...))
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to get alerts")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to get alerts", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(alerts)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to marshal alerts")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to marshal alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) //nolint:errcheck
}

func (a *apiv1) GetAlertsStats(w http.ResponseWriter, r *http.Request, params GetAlertsStatsParams) {
	span := trace.SpanFromContext(r.Context())

	if params.Args == nil {
		params.Args = &map[string]string{}
	}

	q, err := query.UnmarshalAlertStatsQuery(params.Type, *params.Args)
	if err != nil {
		a.logger.Warn().Err(err).Msg("failed to construct alert stats query")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to unmarshal query", http.StatusBadRequest)
		return
	}

	stats, err := a.api.QueryAlertStats(r.Context(), q)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to query alert stats")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to query alert stats", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(stats)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to marshal alert stats")
		span.SetStatus(codes.Error, err.Error())
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
		a.logger.Debug().Err(err).Msg("failed to get cluster status")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if clusterNodes == nil {
		http.Error(w, "no clusterer configured", http.StatusNotFound)
		return
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

func (a *apiv1) PostAlertsAck(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	ack := PostAlertsAckJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&ack); err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, fmt.Sprintf("failed to decode body: %q", err.Error()), http.StatusBadRequest)
		return
	}

	if ack.AlertID == nil {
		span.SetStatus(codes.Error, "missing Alert ID")
		http.Error(w, "missing Alert ID", http.StatusBadRequest)
		return
	}

	alertAck := model.AlertAcknowledgement{
		Creator: ack.Creator,
		Comment: ack.Comment,
	}

	if err := a.api.AckAlert(r.Context(), *ack.AlertID, alertAck); err != nil {
		a.logger.Debug().Err(err).Msg("failed to broadcast alert acknowledgment")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, fmt.Sprintf("failed to handle alert acknowledgment: %q", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *apiv1) PostSilences(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())

	silenceBody := PostSilencesJSONRequestBody{}
	if err := decodeFromContentType(r, &silenceBody); err != nil {
		a.logger.Debug().Err(err).Msg("failed to decode silence")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, fmt.Sprintf("failed to decode silence: %q", err.Error()), http.StatusBadRequest)
		return
	}

	silence := model.Silence{
		Creator:   silenceBody.Creator,
		Comment:   silenceBody.Comment,
		StartTime: silenceBody.StartsAt,
		EndTime:   silenceBody.EndsAt,
		Matchers:  make([]model.Matcher, len(silenceBody.Matchers)),
	}

	for i, matcher := range silenceBody.Matchers {
		silence.Matchers[i] = model.Matcher{
			Label:      matcher.Label,
			Value:      matcher.Value,
			IsRegex:    matcher.IsRegex,
			IsNegative: matcher.IsNegative,
		}
	}

	if err := a.api.PostSilence(r.Context(), silence); err != nil {
		a.logger.Debug().Err(err).Msg("failed to broadcast alert acknowledgment")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, fmt.Sprintf("failed to post silence: %q", err.Error()), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(silence)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to marshal silence")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to marshal silence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes) // nolint:errcheck
}

func (a *apiv1) GetSilences(w http.ResponseWriter, r *http.Request, params GetSilencesParams) {
	span := trace.SpanFromContext(r.Context())

	queries := []query.SilenceFilter{}
	if params.Matchers != nil {
		for _, matcherString := range *params.Matchers {
			matcher := model.Matcher{}
			if err := matcher.UnmarshalText(matcherString); err != nil {
				a.logger.Debug().Err(err).Msgf("failed to unmarshal matcher %q", matcherString)
				span.RecordError(err)
				http.Error(w, fmt.Sprintf("failed to unmarshal matcher: %q", matcherString), http.StatusBadRequest)
				return
			}

			if matcher.Label == "__id__" && !matcher.IsRegex && !matcher.IsNegative {
				// For __id__=value matchers, we can use the ID filter, which is more efficient.
				queries = append(queries, query.ID(matcher.Value))
			} else {
				queries = append(queries, query.Matcher(matcher))
			}
		}
	}

	order := ""
	if params.Order != nil {
		order = string(*params.Order)
	}

	opts, err := constructQueryOpts(params.Limit, params.Offset, params.Sort, order)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to construct query options")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := query.NewSilenceQuery(query.AllSilences(queries...), opts...)

	silences, err := a.api.GetSilences(r.Context(), query)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to get silences")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to get silences", http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(silences)
	if err != nil {
		a.logger.Debug().Err(err).Msg("failed to marshal silences")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to marshal silences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes) // nolint:errcheck // Errors writing here are not recoverable.
}
