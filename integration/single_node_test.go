package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

// Test that Kiora doesn't immediately exit.
func TestKioraStart(t *testing.T) {
	initT(t)
	kiora := NewKioraInstance().Start(t)
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	require.Equal(t, context.DeadlineExceeded, kiora.WaitForExit(ctx), "StdErr: %q", kiora.Stderr())
}

// Test that a post to a Kiora instance stores the alert.
func TestKioraAlertPost(t *testing.T) {
	initT(t)

	kiora := NewKioraInstance().Start(t)

	// Send a bunch of the same alert.
	for i := 0; i < 50; i++ {
		kiora.SendAlert(t, context.TODO(), dummyAlert())
	}

	// Sleep a bit to apply the alert.
	time.Sleep(1 * time.Second)

	require.Contains(t, kiora.stdout.String(), "foo")

	// It should only have fired once.
	require.Equal(t, 1, strings.Count(kiora.stdout.String(), "foo"))
}

// Test that an alert refires if it fires, resolves, and then refires.
func TestKioraResolveResends(t *testing.T) {
	initT(t)

	kiora := NewKioraInstance().Start(t)

	alert := dummyAlert()
	resolved := dummyAlert()
	resolved.Status = model.AlertStatusResolved

	kiora.SendAlert(t, context.Background(), alert)
	time.Sleep(1 * time.Second)
	require.Equal(t, 1, strings.Count(kiora.Stdout(), "foo"))

	kiora.SendAlert(t, context.Background(), resolved)
	time.Sleep(1 * time.Second)
	require.Contains(t, kiora.stdout.String(), "resolved")
	require.Equal(t, 2, strings.Count(kiora.Stdout(), "foo"))

	kiora.SendAlert(t, context.Background(), alert)
	time.Sleep(1 * time.Second)
	require.Equal(t, 3, strings.Count(kiora.Stdout(), "foo"))
}
