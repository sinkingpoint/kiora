package messages

import "github.com/sinkingpoint/kiora/lib/kiora/model"

func init() {
	registerMessage(func() Message { return &Acknowledgement{} })
}

type Acknowledgement struct {
	AlertID         string
	Acknowledgement model.AlertAcknowledgement
}

func (a *Acknowledgement) Name() string {
	return "ack"
}
