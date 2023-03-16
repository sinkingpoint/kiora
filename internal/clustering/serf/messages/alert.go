package messages

import "github.com/sinkingpoint/kiora/lib/kiora/model"

var _ = Message(&AlertMessage{})

type Message interface {
	Name() string
}

type AlertMessage struct {
	Alerts []model.Alert
}

func (a *AlertMessage) Name() string {
	return "alerts"
}
