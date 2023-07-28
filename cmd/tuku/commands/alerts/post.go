package alerts

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/cmd/tuku/commands"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type AlertsPostCmd struct {
	Alerts []string `arg:"" help:"Alerts to post."`
}

func (a *AlertsPostCmd) Run(ctx *commands.Context) error {
	alerts := []model.Alert{}

	for _, alertJSON := range a.Alerts {
		alert := model.Alert{}
		if err := json.Unmarshal([]byte(alertJSON), &alert); err != nil {
			return errors.Wrap(err, "failed to unmarshal alert")
		}

		alerts = append(alerts, alert)
	}

	return ctx.Kiora.PostAlerts(alerts)
}
