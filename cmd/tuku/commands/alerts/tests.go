package alerts

import (
	"math/rand"
	"strconv"

	"github.com/sinkingpoint/kiora/cmd/tuku/commands"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type AlertsTestCmd struct {
	NumPossibleAlerts  int `help:"Number of possible alerts to generate." default:"100"`
	NumAlerts          int `help:"Number of alerts to generate." default:"1000"`
	MaximumLabels      int `help:"Maximum number of labels per alert." default:"10"`
	MaximumCardinality int `help:"Maximum cardinality of each label." default:"100"`
}

type PotentialAlerts struct {
	Name               string
	LabelNames         []string
	LabelCardinalities []int
}

func (a *AlertsTestCmd) generatePotentialAlerts() []PotentialAlerts {
	potentialAlerts := make([]PotentialAlerts, 100)

	// Generate 100 potential alerts with random labels and cardinalities.
	for i := 0; i < a.NumPossibleAlerts; i++ {
		numLabels := rand.Intn(a.MaximumLabels) + 1
		labelNames := make([]string, numLabels)
		labelCardinalities := make([]int, numLabels)
		for j := 0; j < numLabels; j++ {
			labelNames[j] = "Label_" + strconv.Itoa(rand.Intn(a.MaximumLabels)+1)
			labelCardinalities[j] = rand.Intn(a.MaximumCardinality) + 1
		}

		potentialAlerts[i] = PotentialAlerts{
			Name:               "Alert_" + strconv.Itoa(i+1),
			LabelNames:         labelNames,
			LabelCardinalities: labelCardinalities,
		}
	}

	return potentialAlerts
}

func (a *AlertsTestCmd) generateAlerts() []model.Alert {
	// Generate 1000 alerts based on a random template from PotentialAlerts
	alerts := make([]model.Alert, a.NumAlerts)
	potentialAlerts := a.generatePotentialAlerts()
	for i := 0; i < len(alerts); i++ {
		potentialIndex := rand.Intn(len(potentialAlerts))
		potential := potentialAlerts[potentialIndex]
		alert := model.Alert{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Status:      model.AlertStatusFiring,
		}

		alert.Labels["alertname"] = potential.Name

		for j := 0; j < len(potential.LabelNames); j++ {
			labelName := potential.LabelNames[j]
			labelCardinality := potential.LabelCardinalities[j]
			k := rand.Intn(labelCardinality)
			labelValue := strconv.Itoa(k)
			alert.Labels[labelName] = labelValue
		}

		alerts[i] = alert
	}

	return alerts
}

func (a *AlertsTestCmd) Run(ctx *commands.Context) error {
	alerts := a.generateAlerts()

	return ctx.Kiora.PostAlerts(alerts)
}
