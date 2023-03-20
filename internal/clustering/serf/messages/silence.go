package messages

import "github.com/sinkingpoint/kiora/lib/kiora/model"

func init() {
	registerMessage(func() Message { return &Silence{} })
}

type Silence struct {
	Silence model.Silence
}

func (s *Silence) Name() string {
	return "silence"
}
