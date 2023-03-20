package alerts

import (
	"fmt"

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
		return err
	}

	fmt.Println(string(out))

	return nil
}
