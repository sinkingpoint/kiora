package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoad(t *testing.T) {
	tests := []struct {
		name          string
		config        string
		expectSuccess bool
	}{
		{
			name: "standard config",
			config: `digraph Config {
				console_debug [type="stdout"];
				alerts -> console_debug;
			}`,
			expectSuccess: true,
		},
		{
			name: "cycle config",
			config: `digraph Config {
				console_debug [type="stdout"];
				alerts -> console_debug -> alerts;
			}`,
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "kiora_test")
			require.NoError(t, err)
			_, err = file.Write([]byte(tt.config))
			require.NoError(t, err)
			file.Close()

			require.NoError(t, err)
			_, err = config.LoadConfigFile(file.Name())
			if tt.expectSuccess {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestConfigAckFilter(t *testing.T) {
	tests := []struct {
		name          string
		config        string
		alert         *model.Alert
		ack           *model.AlertAcknowledgement
		expectSuccess bool
	}{
		{
			name: "noop config",
			config: `digraph Config {
			}`,
			alert:         &model.Alert{},
			ack:           &model.AlertAcknowledgement{},
			expectSuccess: true,
		},
		{
			name: "bad email",
			config: `digraph Config {
				email_filter -> acks [type="regex" field="from" regex=".*@example.com"];
			}`,
			alert: &model.Alert{},
			ack: &model.AlertAcknowledgement{
				By: "colin@notanemail",
			},
			expectSuccess: false,
		},
		{
			name: "good email",
			config: `digraph Config {
				email_filter -> acks [type="regex" field="from" regex=".*@example.com"];
			}`,
			alert: &model.Alert{},
			ack: &model.AlertAcknowledgement{
				By: "colin@example.com",
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "kiora_test")
			require.NoError(t, err)
			_, err = file.Write([]byte(tt.config))
			require.NoError(t, err)
			file.Close()

			require.NoError(t, err)
			cfg, err := config.LoadConfigFile(file.Name())
			require.NoError(t, err)

			acceptable := cfg.AlertAcknowledgementIsValid(context.TODO(), tt.ack)
			assert.Equal(t, tt.expectSuccess, acceptable)
		})
	}
}
