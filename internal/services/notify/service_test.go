package notify

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/kiora/mocks/mock_clustering"
	"github.com/sinkingpoint/kiora/mocks/mock_config"
	"github.com/sinkingpoint/kiora/mocks/mock_kioradb"
	"github.com/sinkingpoint/kiora/mocks/mock_services"
)

func TestNotifyServiceFiring(t *testing.T) {
	type test struct {
		Name              string
		Alerts            []model.Alert
		ExpectedBroadcast []int
	}

	tests := []test{
		{
			Name: "test_notify_fires",
			Alerts: []model.Alert{
				{
					Status:         model.AlertStatusFiring,
					LastNotifyTime: time.Time{},
				},
			},
			ExpectedBroadcast: []int{0},
		},
		{
			Name: "test_resolved_fires",
			Alerts: []model.Alert{
				// An alert that has passed its EndTime, but its LastNotifyTime is before the EndTime (i.e. it's resolved, but we haven't send the resolve)
				{
					Status:         model.AlertStatusResolved,
					EndTime:        time.Time{}.Add(1 * time.Hour),
					LastNotifyTime: time.Time{},
				},
			},
			ExpectedBroadcast: []int{0},
		},
		{
			Name: "test_resolved_doesnt_refire",
			Alerts: []model.Alert{
				{
					Status:         model.AlertStatusResolved,
					EndTime:        time.Time{},
					LastNotifyTime: time.Time{}.Add(1 * time.Hour),
				},
			},
			ExpectedBroadcast: []int{},
		},
	}

	testTime := time.Now()
	stubs.Time.Now = func() time.Time {
		return testTime
	}

	ctrl := gomock.NewController(t)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			bus := mock_services.NewMockBus(ctrl)
			bus.EXPECT().DB().DoAndReturn(func() kioradb.DB {
				return mock_kioradb.MockDBWithAlerts(ctrl, tt.Alerts)
			}).AnyTimes()

			alerts := []model.Alert{}
			for _, idx := range tt.ExpectedBroadcast {
				alert := tt.Alerts[idx]
				alert.LastNotifyTime = testTime
				alerts = append(alerts, alert)
			}

			bus.EXPECT().Broadcaster().Return(mock_clustering.MockBroadcasterExpectingAlerts(ctrl, alerts)).AnyTimes()

			notifier := mock_config.NewMockNotifier(ctrl)
			notifier.EXPECT().Notify(gomock.Any(), alerts).AnyTimes()

			conf := mock_config.NewMockConfig(ctrl)
			conf.EXPECT().GetNotifiersForAlert(gomock.Any(), gomock.Any()).Return([]config.NotifierSettings{config.NewNotifier(notifier).WithGroupWait(0)}).AnyTimes()

			notifyService := NewNotifyService(conf, bus)
			notifyService.notifyFiring(context.TODO())
			notifyService.notifyResolved(context.TODO())
		})
	}
}
