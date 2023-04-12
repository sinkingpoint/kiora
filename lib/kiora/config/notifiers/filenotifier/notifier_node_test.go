package filenotifier_test

import (
	"context"
	"os"
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/config/notifiers/filenotifier"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileNotifierNode(t *testing.T) {
	file, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer file.Close()

	node, err := filenotifier.New("",
		map[string]string{
			"type": "file",
			"path": file.Name(),
		})

	require.NoError(t, err)

	processor := node.(config.Notifier)

	assert.NoError(t, processor.Notify(context.Background(), model.Alert{
		Labels: model.Labels{
			"alertname": "foo",
		},
	}))

	fileContents, err := os.ReadFile(file.Name())
	require.NoError(t, err)

	assert.Contains(t, string(fileContents), "alertname")
	assert.Contains(t, string(fileContents), "foo")
}
