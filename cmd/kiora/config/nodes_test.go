package config

import (
	"context"
	"os"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileNotifyNode(t *testing.T) {
	file, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer file.Close()

	node, err := NewFileNotifierNode(node{
		name: "",
		attrs: map[string]string{
			"type": "file",
			"path": file.Name(),
		},
	})

	require.NoError(t, err)

	processor := node.(kioradb.ModelWriter)

	assert.NoError(t, processor.ProcessAlerts(context.Background(), model.Alert{
		Labels: model.Labels{
			"alertname": "foo",
		},
	}))

	fileContents, err := os.ReadFile(file.Name())
	require.NoError(t, err)

	assert.Contains(t, string(fileContents), "alertname")
	assert.Contains(t, string(fileContents), "foo")
}
