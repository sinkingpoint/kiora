package alerts

import (
	"time"

	"github.com/sinkingpoint/kiora/cmd/tuku/commands"
	"github.com/sinkingpoint/kiora/internal/testutils"
)

type AlertsTestCmd struct {
	NumPossibleAlerts  int `help:"Number of possible alerts to generate." default:"100"`
	NumAlerts          int `help:"Number of alerts to generate." default:"1000"`
	MaximumLabels      int `help:"Maximum number of labels per alert." default:"10"`
	MaximumCardinality int `help:"Maximum cardinality of each label." default:"100"`
	BatchSize          int `help:"Number of alerts to send in each batch." default:"100"`
}

func (a *AlertsTestCmd) Run(ctx *commands.Context) error {
	alerts := testutils.GenerateDummyAlerts(a.NumAlerts, a.NumPossibleAlerts, a.MaximumLabels, a.MaximumCardinality)

	startTime := time.Now()

	for i := 0; i < len(alerts); i += a.BatchSize {
		end := i + a.BatchSize
		if end > len(alerts) {
			end = len(alerts)
		}

		for j := i; j < end; j++ {
			alerts[j].StartTime = startTime.Add(time.Duration(i+j) * time.Second)
		}

		if err := ctx.Kiora.PostAlerts(alerts[i:end]); err != nil {
			return err
		}
	}

	return ctx.Kiora.PostAlerts(alerts)
}
