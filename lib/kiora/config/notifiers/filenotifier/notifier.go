package filenotifier

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/sinkingpoint/kiora/internal/encoding"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"go.opentelemetry.io/otel"
)

func init() {
	config.RegisterNode(STDOUT_NODE_NAME, New)
	config.RegisterNode(STDERR_NODE_NAME, New)
	config.RegisterNode(FILE_NODE_NAME, New)
}

const STDOUT_NODE_NAME = "stdout"
const STDERR_NODE_NAME = "stderr"
const FILE_NODE_NAME = "file"
const DEFAULT_ENCODING = "json"

var _ = config.Notifier(&FileNotifier{})

// FileNotifier represents a node that can output alerts to a Writer.
type FileNotifier struct {
	name    config.NotifierName
	encoder encoding.Encoder
	file    io.WriteCloser
}

func New(name string, bus config.NodeBus, attrs map[string]string) (config.Node, error) {
	encodingName := DEFAULT_ENCODING
	if enc, ok := attrs["encoding"]; ok {
		encodingName = enc
	}

	encoder := encoding.LookupEncoding(encodingName)
	if encoder == nil {
		return nil, fmt.Errorf("invalid encoding: %q", encodingName)
	}

	switch attrs["type"] {
	case "stdout":
		return &FileNotifier{
			name:    config.NotifierName(name),
			encoder: encoder,
			file:    os.Stdout,
		}, nil
	case "stderr":
		return &FileNotifier{
			name:    config.NotifierName(name),
			encoder: encoder,
			file:    os.Stderr,
		}, nil
	case "", "file":
		fileName := attrs["path"]
		if fileName == "" {
			return nil, errors.New("missing `path` in file node")
		}

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %q in file node: %w", fileName, err)
		}

		return &FileNotifier{
			name:    config.NotifierName(name),
			encoder: encoder,
			file:    file,
		}, nil
	default:
		return nil, fmt.Errorf("invalid type for file node: %q", attrs["type"])
	}
}

func (f *FileNotifier) Name() config.NotifierName {
	return f.name
}

func (f *FileNotifier) Type() string {
	return "file"
}

func (f *FileNotifier) Notify(ctx context.Context, alerts ...model.Alert) *config.NotificationError {
	_, span := otel.Tracer("").Start(ctx, "FileNotifierNode.Notify")
	defer span.End()

	var lastError error
	for _, alert := range alerts {
		bytes, err := f.encoder.Marshal(alert)
		if err != nil {
			lastError = multierror.Append(lastError, err)
			continue
		}

		bytes = append(bytes, '\n')

		if _, err := f.file.Write(bytes); err != nil {
			lastError = multierror.Append(lastError, err)
		}
	}

	if lastError != nil {
		return config.NewNotificationError(lastError, false)
	}

	return nil
}
