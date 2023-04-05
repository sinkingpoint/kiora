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

// TestNotifyServiceNotifies tests that the notify service will send notifications for alerts that are firing or resolved.
func TestNotifyServiceNotifies(t *testing.T) {
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
			alerts := []model.Alert{}
			for _, idx := range tt.ExpectedBroadcast {
				alert := tt.Alerts[idx]
				alert.LastNotifyTime = testTime
				alerts = append(alerts, alert)
			}

			bus := mock_services.NewMockBus(ctrl)
			bus.EXPECT().DB().DoAndReturn(func() kioradb.DB {
				return mock_kioradb.MockDBWithAlerts(ctrl, tt.Alerts)
			}).MinTimes(1)

			notifier := mock_config.NewMockNotifier(ctrl)

			conf := mock_config.NewMockConfig(ctrl)

			if len(tt.ExpectedBroadcast) > 0 {
				bus.EXPECT().Broadcaster().Return(mock_clustering.MockBroadcasterExpectingAlerts(ctrl, alerts)).MinTimes(1)
				notifier.EXPECT().Notify(gomock.Any(), alerts).Times(1)

				for _, i := range tt.ExpectedBroadcast {
					conf.EXPECT().GetNotifiersForAlert(gomock.Any(), &tt.Alerts[i]).Return([]config.NotifierSettings{
						config.NewNotifier(notifier).WithGroupWait(0),
					}).Times(1)
				}
			}

			notifyService := NewNotifyService(conf, bus)
			notifyService.notifyFiring(context.TODO())
			notifyService.notifyResolved(context.TODO())
			notifyService.notifyGroup(context.TODO())
		})
	}
}

// TestNotifyServiceGrouping tests that the notify service will group alerts by their labels.
func TestNotifyServiceGrouping(t *testing.T) {
	rawAlerts := []model.Alert{
		{
			Labels: model.Labels{
				"foo": "bar",
				"bar": "baz",
			},
			Annotations: map[string]string{},
			Status:      model.AlertStatusFiring,
			StartTime:   stubs.Time.Now(),
		},
		{
			Labels: model.Labels{
				"foo": "bar",
				"bar": "foo",
			},
			Annotations: map[string]string{},
			Status:      model.AlertStatusFiring,
			StartTime:   stubs.Time.Now(),
		},
		{
			Labels: model.Labels{
				"foo": "baz",
				"bar": "foo",
			},
			Annotations: map[string]string{},
			Status:      model.AlertStatusFiring,
			StartTime:   stubs.Time.Now(),
		},
	}

	testTime := time.Now()
	stubs.Time.Now = func() time.Time {
		return testTime
	}

	expectedGroups := [][]int{{0, 1}, {2}}

	ctrl := gomock.NewController(t)
	notifier := mock_config.NewMockNotifier(ctrl)
	notifier.EXPECT().Name().Return(config.NotifierName("mock notifier")).MinTimes(1)

	groupedAlerts := [][]model.Alert{}
	for _, group := range expectedGroups {
		alertGroup := []model.Alert{}
		for _, idx := range group {
			alert := rawAlerts[idx]
			alert.LastNotifyTime = testTime
			alertGroup = append(alertGroup, alert)
		}

		groupedAlerts = append(groupedAlerts, alertGroup)
	}

	bus := mock_services.NewMockBus(ctrl)
	bus.EXPECT().DB().DoAndReturn(func() kioradb.DB {
		return mock_kioradb.MockDBWithAlerts(ctrl, rawAlerts)
	}).MinTimes(1)

	for _, group := range groupedAlerts {
		notifier.EXPECT().Notify(gomock.Any(), group).Times(1)
	}

	conf := mock_config.NewMockConfig(ctrl)

	if len(expectedGroups) > 0 {
		// Expect that we'll get a single broadcast call for each alert group.
		bus.EXPECT().Broadcaster().Return(mock_clustering.MockBroadcasterExpectingAlerts(ctrl, groupedAlerts...)).MinTimes(1)

		// Expect that we'll get a single call to the notifier config for each alert group.
		notifierConf := config.NewNotifier(notifier).WithGroupWait(time.Millisecond).WithGroupLabels("foo")

		conf.EXPECT().GetNotifiersForAlert(gomock.Any(), gomock.Any()).Return([]config.NotifierSettings{notifierConf}).Times(len(rawAlerts))
	}

	notifyService := NewNotifyService(conf, bus)
	notifyService.notifyFiring(context.TODO())
	notifyService.notifyResolved(context.TODO())

	// Wait for the group to time out.
	time.Sleep(10 * time.Millisecond)
	notifyService.notifyGroup(context.TODO())
}
