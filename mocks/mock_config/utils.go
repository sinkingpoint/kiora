package mock_config

import gomock "github.com/golang/mock/gomock"

// NewMockConfigAllowingEverything returns a MockConfig set up with a `ValidateData` method that returns nil.
func NewMockConfigAllowingEverything(ctrl *gomock.Controller) *MockConfig {
	conf := NewMockConfig(ctrl)
	conf.EXPECT().ValidateData(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	return conf
}
