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

			db := mock_kioradb.MockDBWithAlerts(ctrl, tt.Alerts)

			if len(tt.ExpectedBroadcast) > 0 {
				db.EXPECT().StoreAlerts(gomock.Any(), alerts).Times(1)
			}

			bus := mock_services.NewMockBus(ctrl)
			bus.EXPECT().DB().DoAndReturn(func() kioradb.DB {
				return db
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
	testTime := time.Now()
	stubs.Time.Now = func() time.Time {
		return testTime
	}

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

	ctrl := gomock.NewController(t)
	db := mock_kioradb.MockDBWithAlerts(ctrl, rawAlerts)
	broadcaster := mock_clustering.NewMockBroadcaster(ctrl)
	notifier := mock_config.NewMockNotifier(ctrl)
	notifier.EXPECT().Name().Return(config.NotifierName("mock notifier")).MinTimes(1)

	// The first two alerts have the same `foo` label, so they should be grouped together. The third is different.
	expectedGroups := [][]int{{0, 1}, {2}}
	for _, groupIndexes := range expectedGroups {
		// For each group, we expect the notifier to be called to notify the group, and the broadcaster to be called to broadcast the group.
		group := []interface{}{}
		for _, idx := range groupIndexes {
			alert := rawAlerts[idx]
			alert.LastNotifyTime = testTime
			group = append(group, alert)
		}

		notifier.EXPECT().Notify(gomock.Any(), group...).Times(1)
		broadcaster.EXPECT().BroadcastAlerts(gomock.Any(), group...).Times(1)
	}

	bus := mock_services.NewMockBus(ctrl)
	bus.EXPECT().DB().Return(db).MinTimes(1)
	bus.EXPECT().Broadcaster().Return(broadcaster).MinTimes(1)

	// For every alert, we expect the config to be queried for notifiers for that alert.
	conf := mock_config.NewMockConfig(ctrl)
	for i := range rawAlerts {
		alert := rawAlerts[i]
		conf.EXPECT().GetNotifiersForAlert(gomock.Any(), &alert).Return([]config.NotifierSettings{
			config.NewNotifier(notifier).WithGroupWait(1 * time.Millisecond).WithGroupLabels("foo"),
		}).Times(1)
	}

	// For every alert, we expect the alert to be stored and broadcast after the notifyFiring call, but _before_ the notifyGroup call.
	for i := range rawAlerts {
		alert := rawAlerts[i]
		alert.LastNotifyTime = stubs.Time.Now()
		db.EXPECT().StoreAlerts(gomock.Any(), []model.Alert{alert}).Times(1)
		broadcaster.EXPECT().BroadcastAlerts(gomock.Any(), []model.Alert{alert}).Times(1)
	}

	notifyService := NewNotifyService(conf, bus)
	notifyService.notifyFiring(context.TODO())
	notifyService.notifyResolved(context.TODO())

	// Wait 10ms for the group to expire.
	time.Sleep(10 * time.Millisecond)
	notifyService.notifyGroup(context.TODO())
}
