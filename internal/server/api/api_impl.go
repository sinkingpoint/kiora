package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ = API(&APIImpl{})

// API defines an interface that represents all the operations that can be performed on the kiora API.
type API interface {
	// GetAlerts returns a list of alerts matching the given query.
	GetAlerts(ctx context.Context, q *query.AlertQuery) ([]model.Alert, error)

	// PostAlerts stores the given alerts in the database, updating any existing alerts with the same labels.
	PostAlerts(ctx context.Context, alerts []model.Alert) error

	// QueryAlertStats executes the given stats query, returning the resulting frames.
	QueryAlertStats(ctx context.Context, q query.AlertStatsQuery) ([]query.StatsResult, error)

	// GetSilences returns a list of silences matching the given query.
	GetSilences(ctx context.Context) ([]model.Silence, error)

	// PostSilences stores the given silences in the database, updating any existing silences with the same ID.
	PostSilence(ctx context.Context, silences model.Silence) error

	// AckAlert acknowledges the given alert with the given acknowledgement.
	AckAlert(ctx context.Context, alertID string, alertAck model.AlertAcknowledgement) error

	// GetClusterStatus returns the status of the nodes in the cluster.
	GetClusterStatus(ctx context.Context) ([]any, error)
}

type APIImpl struct {
	bus       services.Bus
	clusterer clustering.Clusterer
}

func NewAPIImpl(bus services.Bus, clusterer clustering.Clusterer) *APIImpl {
	return &APIImpl{
		bus:       bus,
		clusterer: clusterer,
	}
}

func (a *APIImpl) GetAlerts(ctx context.Context, q *query.AlertQuery) ([]model.Alert, error) {
	return a.bus.DB().QueryAlerts(ctx, q), nil
}

func (a *APIImpl) PostAlerts(ctx context.Context, alerts []model.Alert) error {
	return a.bus.Broadcaster().BroadcastAlerts(ctx, alerts...)
}

func (a *APIImpl) QueryAlertStats(ctx context.Context, q query.AlertStatsQuery) ([]query.StatsResult, error) {
	return kioradb.QueryAlertStats(ctx, a.bus.DB(), q)
}

func (a *APIImpl) GetSilences(ctx context.Context) ([]model.Silence, error) {
	return a.bus.DB().QuerySilences(ctx, query.MatchAll()), nil
}

func (a *APIImpl) PostSilence(ctx context.Context, silence model.Silence) error {
	if err := a.bus.Config().ValidateData(ctx, &silence); err != nil {
		return err
	}

	return a.bus.Broadcaster().BroadcastSilences(ctx, silence)
}

func (a *APIImpl) AckAlert(ctx context.Context, alertID string, alertAck model.AlertAcknowledgement) error {
	if err := a.bus.Config().ValidateData(ctx, &alertAck); err != nil {
		return err
	}

	if len(a.bus.DB().QueryAlerts(ctx, query.NewAlertQuery(query.ID(alertID)))) == 0 {
		return fmt.Errorf("alert %q not found", alertID)
	}

	return a.bus.Broadcaster().BroadcastAlertAcknowledgement(ctx, alertID, alertAck)
}

func (a *APIImpl) GetClusterStatus(ctx context.Context) ([]any, error) {
	if a.clusterer == nil {
		return nil, errors.New("no clusterer configured")
	}

	return a.clusterer.Nodes(), nil
}
