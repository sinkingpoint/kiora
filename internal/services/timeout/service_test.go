package timeout

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/kiora/mocks/mock_clustering"
	"github.com/sinkingpoint/kiora/mocks/mock_kioradb"
	"github.com/sinkingpoint/kiora/mocks/mock_services"
)

func TestTimeoutService(t *testing.T) {
	type test struct {
		Name              string
		Alerts            []model.Alert
		ExpectedBroadcast []int
	}

	tests := []test{
		{
			Name: "test_no_timed_out",
			Alerts: []model.Alert{
				{
					TimeOutDeadline: time.Now().Add(1 * time.Hour),
					Status:          model.AlertStatusFiring,
				},
			},
			ExpectedBroadcast: []int{},
		},
		{
			Name: "test_timed_out",
			Alerts: []model.Alert{
				{
					TimeOutDeadline: time.Now().Add(-1 * time.Hour),
					Status:          model.AlertStatusFiring,
				},
			},
			ExpectedBroadcast: []int{0},
		},

		{
			Name: "test_resolved_doesn't_time_out",
			Alerts: []model.Alert{
				{
					TimeOutDeadline: time.Now().Add(-1 * time.Hour),
					Status:          model.AlertStatusResolved,
				},
			},
			ExpectedBroadcast: []int{},
		},
	}

	ctrl := gomock.NewController(t)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			bus := mock_services.NewMockBus(ctrl)
			bus.EXPECT().DB().Return(mock_kioradb.MockDBWithAlerts(ctrl, tt.Alerts))

			expectedBroadcast := []model.Alert{}
			for _, idx := range tt.ExpectedBroadcast {
				alert := tt.Alerts[idx]
				alert.Status = model.AlertStatusTimedOut
				expectedBroadcast = append(expectedBroadcast, alert)
			}
			bus.EXPECT().Broadcaster().Return(mock_clustering.MockBroadcasterExpectingAlerts(ctrl, expectedBroadcast))

			timeoutService := NewTimeoutService(bus)
			timeoutService.timeoutAlerts(context.TODO())
		})
	}
}
