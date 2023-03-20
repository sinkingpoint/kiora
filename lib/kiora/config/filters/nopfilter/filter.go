package nopfilter

import (
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

func init() {
	config.RegisterFilter("", NewNopFilter)
}

// NopFilter is the default filter when a type is not specified. It does nothing and lets everything through.
type NopFilter struct{}

func NewNopFilter(attrs map[string]string) (config.Filter, error) {
	return &NopFilter{}, nil
}

func (n *NopFilter) Type() string {
	return "nop"
}

func (n *NopFilter) FilterAlert(a *model.Alert) bool {
	return true
}

func (n *NopFilter) FilterAlertAcknowledgement(alert *model.Alert, ack *model.AlertAcknowledgement) bool {
	return true
}
