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

func writeConfigFile(t *testing.T, config string) string {
	t.Helper()
	file, err := os.CreateTemp("", "kiora_test")
	require.NoError(t, err)
	_, err = file.Write([]byte(config))
	require.NoError(t, err)
	file.Close()
	return file.Name()
}

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
			fileName := writeConfigFile(t, tt.config)
			_, err := config.LoadConfigFile(fileName)
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
		name        string
		config      string
		ack         *model.AlertAcknowledgement
		expectError []string
	}{
		{
			name: "noop config",
			config: `digraph Config {
			}`,
			ack:         &model.AlertAcknowledgement{},
			expectError: nil,
		},
		{
			name: "bad email",
			config: `digraph Config {
				email_filter -> acks [type="regex" field="from" regex=".*@example.com"];
			}`,
			ack: &model.AlertAcknowledgement{
				By: "colin@notanemail",
			},
			expectError: []string{
				"field from doesn't match",
			},
		},
		{
			name: "good email",
			config: `digraph Config {
				email_filter -> acks [type="regex" field="from" regex=".*@example.com"];
			}`,
			ack: &model.AlertAcknowledgement{
				By: "colin@example.com",
			},
			expectError: nil,
		},
		{
			name: "two step validation",
			config: `digraph config {
				console [type="stdout"];
				alerts -> console;

				test_email -> test_comment [type="regex" field="from" regex=".+@example.com"];
				test_comment -> acks [type="regex" field="comment" regex=".+"];
			}`,
			ack: &model.AlertAcknowledgement{
				By: "colin@example.com",
			},
			expectError: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := writeConfigFile(t, tt.config)
			cfg, err := config.LoadConfigFile(fileName)
			require.NoError(t, err)

			err = cfg.AlertAcknowledgementIsValid(context.TODO(), tt.ack)
			if tt.expectError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectError != nil && err == nil {
				t.Fatal("expected error, got none")
			}

			for _, s := range tt.expectError {
				assert.Contains(t, err.Error(), s, "expected error to contain %q", s)
			}
		})
	}
}
