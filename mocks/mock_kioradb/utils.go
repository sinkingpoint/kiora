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
		db.EXPECT().QueryAlerts(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, query query.AlertFilter) []model.Alert {
			ret := []model.Alert{}
			for _, alert := range alerts {
				if query.MatchesAlert(ctx, &alert) {
					ret = append(ret, alert)
				}
			}

			return ret
		}).Times(1)
	}

	return db
}
