package config

import (
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/filters/duration"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/filters/nop"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/filters/regex"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/filenotifier"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/slack"
)
