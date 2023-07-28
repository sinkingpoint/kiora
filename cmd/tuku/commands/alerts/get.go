package alerts

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/cmd/tuku/commands"
)

type AlertsGetCmd struct{}

func (a *AlertsGetCmd) Run(ctx *commands.Context) error {
	alerts, err := ctx.Kiora.GetAlerts()
	if err != nil {
		return err
	}

	out, err := ctx.Formatter.Marshal(alerts)
	if err != nil {
		return errors.Wrap(err, "failed to marshal alerts")
	}

	fmt.Println(string(out))

	return nil
}
