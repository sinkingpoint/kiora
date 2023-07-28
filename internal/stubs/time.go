package stubs

import "time"

// stubTime provides a stub over the `time` package so we can control return values in tests.
type stubTime struct {
	Now func() time.Time
}

var Time = stubTime{
	Now: time.Now,
}
