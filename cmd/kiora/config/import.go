package config

import (
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/filters/durationfilter"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/filters/nopfilter"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/filters/regexfilter"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/filenotifier"
	_ "github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/slack"
)
