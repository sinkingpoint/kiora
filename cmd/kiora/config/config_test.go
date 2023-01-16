package config_test

import (
	"os"
	"testing"

	"github.com/sinkingpoint/kiora/cmd/kiora/config"
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
