package mock_kioradb

import (
	context "context"

	"github.com/golang/mock/gomock"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// MockDBWithAlerts returns a MockDB that expects a number of calls to QueryAlerts. The returned
// values are the result of the given query applied to the given alerts.
func MockDBWithAlerts(ctrl *gomock.Controller, alerts ...[]model.Alert) *MockDB {
	db := NewMockDB(ctrl)
	for _, alerts := range alerts {
		db.EXPECT().QueryAlerts(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, query query.AlertQuery) []model.Alert {
			ret := []model.Alert{}
			for _, alert := range alerts {
				if query.Filter.MatchesAlert(ctx, &alert) {
					ret = append(ret, alert)
				}
			}

			return ret
		}).MinTimes(1)
	}

	return db
}
