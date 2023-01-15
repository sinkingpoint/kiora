package kiora

import (
	"context"
	"fmt"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

// NotifierProcessor is responsible for
type NotifierProcessor struct {
	me string
}

func NewNotifierProcessor(myName string) *NotifierProcessor {
	return &NotifierProcessor{
		me: myName,
	}
}

func (s *NotifierProcessor) ProcessAlert(ctx context.Context, broadcast kioradb.ModelWriter, db kioradb.DB, existingAlert, newAlert *model.Alert) error {
	if newAlert.AuthNode != s.me && (existingAlert == nil || existingAlert.AuthNode != s.me) {
		return nil
	}

	if newAlert.Status != model.AlertStatusProcessing {
		return nil
	}

	fmt.Printf("%s notifying for %q\n", s.me, newAlert)

	newAlert.Status = model.AlertStatusFiring

	// TODO(cdouch): Actually fire the alert.

	return broadcast.ProcessAlerts(ctx, *newAlert)
}
