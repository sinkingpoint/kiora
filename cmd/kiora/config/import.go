package config

import (
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/duration"
	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/nop"
	"github.com/sinkingpoint/kiora/lib/kiora/config/filters/regex"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/filenotifier"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/slack"
)

func RegisterNodes() {
	config.RegisterFilter("", nop.NewFilter)
	config.RegisterFilter("regex", regex.NewFilter)
	config.RegisterFilter("duration", duration.NewFilter)
}
