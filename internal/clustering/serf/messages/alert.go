package messages

import "github.com/sinkingpoint/kiora/lib/kiora/model"

var _ = Message(&Alert{})

func init() {
	registerMessage(func() Message { return &Alert{} })
}

// Alert is a message representing an update to an alert.
type Alert struct {
	Alert model.Alert
}

func (a *Alert) Name() string {
	return "alert"
}
