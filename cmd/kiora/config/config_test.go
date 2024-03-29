package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/require"
)

func writeConfigFile(t *testing.T, config string) string {
	t.Helper()
	file, err := os.CreateTemp("", "kiora_test")
	require.NoError(t, err)
	_, err = file.WriteString(config)
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
		{
			name: "tenant key",
			config: `digraph Config {
				tenant_key = "foo-{{ .tenant }}"
			}`,
			expectSuccess: true,
		},
		{
			name: "unknown global",
			config: `digraph Config {
				foo = "bar"
			}`,
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.RegisterNodes()
			fileName := writeConfigFile(t, tt.config)
			_, err := config.LoadConfigFile(fileName, zerolog.New(os.Stdout))
			if tt.expectSuccess {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestConfigAckFilter(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		ack         *model.AlertAcknowledgement
		expectError bool
	}{
		{
			name: "noop config",
			config: `digraph Config {
			}`,
			ack:         &model.AlertAcknowledgement{},
			expectError: false,
		},
		{
			name: "bad email",
			config: `digraph Config {
				email_filter -> acks [type="regex" field="__creator__" regex=".*@example.com"];
			}`,
			ack: &model.AlertAcknowledgement{
				Creator: "colin@notanemail",
			},
			expectError: true,
		},
		{
			name: "good email",
			config: `digraph Config {
				email_filter -> acks [type="regex" field="__creator__" regex=".*@example.com"];
			}`,
			ack: &model.AlertAcknowledgement{
				Creator: "colin@example.com",
			},
			expectError: false,
		},
		{
			name: "two step validation",
			config: `digraph config {
				console [type="stdout"];
				alerts -> console;

				test_email -> test_comment [type="regex" field="__creator__" regex=".+@example.com"];
				test_comment -> acks [type="regex" field="__comment__" regex=".+"];
			}`,
			ack: &model.AlertAcknowledgement{
				Creator: "colin@example.com",
			},
			expectError: true,
		},
		{
			name: "multiple paths",
			config: `digraph config {
				// Sometimes it's useful to have multiple potential validation paths. For example, we might have a bot account
				// that should also be allowed to acknowledge alerts. To do this, we can specify multiple paths into the acks pseudonode.
			
				// First, the regular human path, which must have an email and a comment.
				test_email -> test_comment [type="regex" field="__creator__" regex=".+@example.com"]; // First check the email
				test_comment -> acks [type="regex" field="__comment__" regex=".+"]; // Then check the comment.
			
				// And then a bot path where we don't need a comment, if the from is RespectTables:
				test_respect_tables -> acks [type="regex" field="__creator__" regex="RespectTables"];
			}`,
			ack: &model.AlertAcknowledgement{
				Creator: "colin@example.com",
				Comment: "I'm a human, I promise!",
			},
			expectError: false,
		},
		{
			name: "multiple paths 2",
			config: `digraph config {
				// Sometimes it's useful to have multiple potential validation paths. For example, we might have a bot account
				// that should also be allowed to acknowledge alerts. To do this, we can specify multiple paths into the acks pseudonode.
			
				// First, the regular human path, which must have an email and a comment.
				test_email -> test_comment [type="regex" field="__creator__" regex=".+@example.com"]; // First check the email
				test_comment -> acks [type="regex" field="__comment__" regex=".+"]; // Then check the comment.
			
				// And then a bot path where we don't need a comment, if the from is RespectTables:
				test_respect_tables -> acks [type="regex" field="__creator__" regex="RespectTables"];
			}`,
			ack: &model.AlertAcknowledgement{
				Creator: "RespectTables",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.RegisterNodes()
			fileName := writeConfigFile(t, tt.config)
			cfg, err := config.LoadConfigFile(fileName, zerolog.New(os.Stdout))
			require.NoError(t, err)

			err = cfg.ValidateData(context.TODO(), tt.ack)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
