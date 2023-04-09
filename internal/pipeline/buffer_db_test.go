package pipeline_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/internal/pipeline"
	"github.com/sinkingpoint/kiora/internal/testutils"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferDBAlerts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Generate 1000 alerts, with 100 possible alert names, 10 labels per alert, and 100 possible values per label.
	alerts := testutils.GenerateDummyAlerts(1000, 100, 10, 100)
	db := kioradb.NewInMemoryDB()

	wg := sync.WaitGroup{}
	buffer := pipeline.NewBufferDB(db, 1000, 1000, 10000, 1*time.Millisecond)
	wg.Add(1)
	go func() {
		require.NoError(t, buffer.Run(ctx))
		wg.Done()
	}()

	require.NoError(t, buffer.StoreAlerts(ctx, alerts...))

	cancel()
	wg.Wait()

	storedAlerts := db.QueryAlerts(context.TODO(), query.NewAlertQuery(query.MatchAll()))
	assert.Len(t, storedAlerts, len(alerts), "stored alerts should match the number of alerts we generated (%d != %d)", len(storedAlerts), len(alerts))
}
