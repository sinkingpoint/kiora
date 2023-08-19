package nop

import (
	"context"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
)

func init() {
	config.RegisterFilter("", NewNopFilter)
}

// NopFilter is the default filter when a type is not specified. It does nothing and lets everything through.
type NopFilter struct{}

func NewNopFilter(attrs map[string]string) (config.Filter, error) {
	return &NopFilter{}, nil
}

func (n *NopFilter) Describe() string {
	panic("BUG: Describe() called on NopFilter which can never reject a model.")
}

func (n *NopFilter) Type() string {
	return "nop"
}

func (n *NopFilter) Filter(ctx context.Context, _ config.Fielder) bool {
	return true
}
