package nodes

import (
	"context"
	"os"
	"testing"

	"github.com/sinkingpoint/kiora/internal/services/notify"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileNotifierNode(t *testing.T) {
	file, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer file.Close()

	node, err := NewFileNotifierNode("",
		map[string]string{
			"type": "file",
			"path": file.Name(),
		})

	require.NoError(t, err)

	processor := node.(notify.Notifier)

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
