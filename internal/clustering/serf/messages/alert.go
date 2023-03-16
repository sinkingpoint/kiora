package messages

import "github.com/sinkingpoint/kiora/lib/kiora/model"

var _ = Message(&AlertMessage{})

var messageRegistry = map[string]func() Message{
	"alerts": func() Message { return &AlertMessage{} },
}

// GetMessage returns a blank message with the given name, or nil if there isn't one registered.
func GetMessage(name string) Message {
	if cons, ok := messageRegistry[name]; ok {
		return cons()
	}

	return nil
}

// Message is a message that can be sent through the Serf gossip channel.
type Message interface {
	Name() string
}

// AlertMessage is a message representing an update to an alert.
type AlertMessage struct {
	Alert model.Alert
}

func (a *AlertMessage) Name() string {
	return "alerts"
}
