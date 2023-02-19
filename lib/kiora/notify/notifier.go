package notify

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type Notifier interface {
	Notify(ctx context.Context, alerts ...model.Alert) error
}
