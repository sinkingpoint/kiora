package testutils

import (
	"math/rand"
	"strconv"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type alertTemplates struct {
	Name               string
	LabelNames         []string
	LabelCardinalities []int
}

func generateAlertTemplates(numPossibleAlerts, maximumLabels, maximumCardinality int) []alertTemplates {
	potentialAlerts := make([]alertTemplates, numPossibleAlerts)

	for i := 0; i < numPossibleAlerts; i++ {
		numLabels := rand.Intn(maximumLabels) + 1
		labelNames := make([]string, numLabels)
		labelCardinalities := make([]int, numLabels)
		for j := 0; j < numLabels; j++ {
			labelNames[j] = "Label_" + strconv.Itoa(rand.Intn(maximumLabels)+1)
			labelCardinalities[j] = rand.Intn(maximumCardinality) + 1
		}

		potentialAlerts[i] = alertTemplates{
			Name:               "Alert_" + strconv.Itoa(i+1),
			LabelNames:         labelNames,
			LabelCardinalities: labelCardinalities,
		}
	}

	return potentialAlerts
}

func GenerateDummyAlerts(num, numPossibleAlerts, maximumLabels, maximumCardinality int) []model.Alert {
	// Generate a bunch of alerts with random labels.
	alerts := make([]model.Alert, num)
	templates := generateAlertTemplates(numPossibleAlerts, maximumLabels, maximumCardinality)

	existingAlerts := make(map[model.LabelsHash]bool)

	for i := 0; i < len(alerts); i++ {
		potentialIndex := rand.Intn(len(templates))
		potential := templates[potentialIndex]
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

		if _, ok := existingAlerts[alert.Labels.Hash()]; ok {
			i--
			continue
		}

		alerts[i] = alert
		existingAlerts[alert.Labels.Hash()] = true
	}

	return alerts
}
