package integration

import (
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// initT sets up a T with the basic settings for an integration test,
// skipping if we're in a short test, and running in parallel.
func initT(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()
}

// dummyAlert returns a basic alert to be used in tests.
func dummyAlert() model.Alert {
	return model.Alert{
		Labels: model.Labels{
			"foo": "bar",
		},
		Annotations: map[string]string{},
		Status:      model.AlertStatusFiring,
		StartTime:   time.Now(),
	}
}
