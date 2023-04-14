package config

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

func init() {
	RegisterNode("group_wait", func(name string, bus NodeBus, attrs map[string]string) (Node, error) {
		rawDuration, ok := attrs["duration"]
		if !ok {
			return nil, errors.New("missing duration attribute for group_wait node")
		}

		duration, err := time.ParseDuration(rawDuration)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse duration in group_wait node")
		}

		return NotifierGroupWait(duration), nil
	})

	RegisterNode("group_labels", func(name string, bus NodeBus, attrs map[string]string) (Node, error) {
		rawLabels, ok := attrs["labels"]
		if !ok {
			return nil, errors.New("missing labels attribute for group_labels node")
		}

		labels := strings.Split(rawLabels, ",")

		return NotifierGroupLabels(labels), nil
	})
}

// NotifierSettingsNode is an interface that can be implemented by config nodes that can be used to configure a NotifierSettings.
type NotifierSettingsNode interface {
	Apply(*NotifierSettings) error
}

// NotifierGroupWait is a NotifierSettingsNode that can be used to set the GroupWait field of a NotifierSettings.
type NotifierGroupWait time.Duration

func (n NotifierGroupWait) Type() string {
	return "group_wait"
}

// Apply sets the GroupWait field of the given NotifierSettings.
func (n NotifierGroupWait) Apply(ns *NotifierSettings) error {
	ns.GroupWait = time.Duration(n)
	return nil
}

// NotifierGroupLabels is a NotifierSettingsNode that can be used to set the GroupLabels field of a NotifierSettings.
type NotifierGroupLabels []string

func (n NotifierGroupLabels) Type() string {
	return "group_labels"
}

// Apply sets the GroupLabels field of the given NotifierSettings.
func (n NotifierGroupLabels) Apply(ns *NotifierSettings) error {
	ns.GroupLabels = n
	return nil
}
