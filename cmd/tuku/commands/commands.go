package commands

import (
	"github.com/sinkingpoint/kiora/cmd/tuku/kiora"
	"github.com/sinkingpoint/kiora/internal/encoding"
)

type Context struct {
	Kiora     *kiora.KioraInstance
	Formatter encoding.Encoder
}
