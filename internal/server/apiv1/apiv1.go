package apiv1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/protobuf/proto"
)

func Register(router *mux.Router, db kioradb.DB) {
	api := apiv1{
		db,
	}
	router.Path("/api/v1/alerts").Methods(http.MethodPost).HandlerFunc(api.postAlerts)
	router.Path("/api/v1/alerts").Methods(http.MethodGet).HandlerFunc(api.getAlerts)
}

// ReadBody reads the body from the request, decoding any Content-Encoding present.
func readBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

type apiv1 struct {
	db kioradb.DB
}

// postAlerts handles the POST /alerts request, decoding a list of alerts
// from the body, and forwarding them to the db.
func (a *apiv1) postAlerts(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var alerts []model.Alert

	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(io.NopCloser(bytes.NewBuffer(body)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&alerts); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode body: %q", err.Error()), http.StatusBadRequest)
			return
		}
	case "application/x-capnp":
		// TODO(cdouch): Move this into a helper function so we're not having to manually decode
		// every struct each time.
		protoAlerts := kioraproto.PostAlertsMessage{}
		if err := proto.Unmarshal(body, &protoAlerts); err != nil {
			http.Error(w, "failed to decode body", http.StatusBadRequest)
			return
		}

		for _, protoAlert := range protoAlerts.Alerts {
			var alert model.Alert
			if err := alert.DeserializeFromProto(protoAlert); err != nil {
				http.Error(w, "failed to decode body", http.StatusBadRequest)
				return
			}
			alerts = append(alerts, alert)
		}
	default:
		http.Error(w, fmt.Sprintf("invalid content-type %q", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
		return
	}

	if err := a.db.ProcessAlerts(r.Context(), alerts...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *apiv1) getAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := a.db.GetAlerts(r.Context())
	if err != nil {
		log.Err(err).Msg("failed to get alerts")
		http.Error(w, "failed to get alerts", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(alerts)
	if err != nil {
		log.Err(err).Msg("failed to get alerts")
		http.Error(w, "failed to get alerts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes) //nolint:errcheck
}
