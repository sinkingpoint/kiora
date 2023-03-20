package query

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// Query is a query that can be run against a DB to pull things out of it.
type SilenceQuery interface {
	MatchesSilence(ctx context.Context, alert *model.Silence) bool
}
