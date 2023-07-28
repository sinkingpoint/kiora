package mock_clustering

import (
	"github.com/golang/mock/gomock"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// MockBroadcasterExpectingAlerts contsructs a MockBroadcaster that expects a number of calls to BroadcastAlerts.
func MockBroadcasterExpectingAlerts(ctrl *gomock.Controller, alerts ...[]model.Alert) *MockBroadcaster {
	broadcaster := NewMockBroadcaster(ctrl)
	for _, alerts := range alerts {
		broadcaster.EXPECT().BroadcastAlerts(gomock.Any(), alerts).Times(1)
	}

	return broadcaster
}
