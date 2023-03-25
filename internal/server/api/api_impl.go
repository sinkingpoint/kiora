package api

import (
	"context"
	"errors"

	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/services"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

var _ = API(&APIImpl{})

type API interface {
	GetAlerts(ctx context.Context, q query.AlertFilter) ([]model.Alert, error)
	PostAlerts(ctx context.Context, alerts []model.Alert) error
	GetSilences(ctx context.Context) ([]model.Silence, error)
	PostSilence(ctx context.Context, silences model.Silence) error
	AckAlert(ctx context.Context, alertID string, alertAck model.AlertAcknowledgement) error
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

func (a *APIImpl) GetAlerts(ctx context.Context, q query.AlertFilter) ([]model.Alert, error) {
	return a.bus.DB().QueryAlerts(ctx, q), nil
}

func (a *APIImpl) PostAlerts(ctx context.Context, alerts []model.Alert) error {
	return a.bus.Broadcaster().BroadcastAlerts(ctx, alerts...)
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

	return a.bus.Broadcaster().BroadcastAlertAcknowledgement(ctx, alertID, alertAck)
}

func (a *APIImpl) GetClusterStatus(ctx context.Context) ([]any, error) {
	if a.clusterer == nil {
		return nil, errors.New("no clusterer configured")
	}

	return a.clusterer.Nodes(), nil
}
