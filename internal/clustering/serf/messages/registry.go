package messages

// Message is a message that can be sent through the Serf gossip channel.
type Message interface {
	Name() string
}

type messageConstructor = func() Message

var messageRegistry = map[string]messageConstructor{}

func registerMessage(cons messageConstructor) {
	messageRegistry[cons().Name()] = cons
}

// GetMessage returns a blank message with the given name, or nil if there isn't one registered.
func GetMessage(name string) Message {
	if cons, ok := messageRegistry[name]; ok {
		return cons()
	}

	return nil
}
