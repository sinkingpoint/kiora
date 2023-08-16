package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestFrontendIsEmbedded tests that the frontend is embedded and can be accessed via the HTTP server.
func TestFrontendIsEmbedded(t *testing.T) {
	initT(t)

	kiora := NewKioraInstance(t).Start()
	time.Sleep(1 * time.Second)

	resp, err := http.Get(kiora.GetHTTPURL("/"))
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
