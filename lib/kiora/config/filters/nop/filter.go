package nop

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

// NopFilter is the default filter when a type is not specified. It does nothing and lets everything through.
type NopFilter struct{}

func NewFilter(globals *config.Globals, attrs map[string]string) (config.Filter, error) {
	return &NopFilter{}, nil
}

func (n *NopFilter) Type() string {
	return "nop"
}

func (n *NopFilter) Filter(ctx context.Context, _ config.Fielder) error {
	return nil
}
