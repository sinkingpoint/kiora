package clustering

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type Broadcaster interface {
	BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error
}
