package config

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/sinkingpoint/kiora/internal/encoding"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

const DEFAULT_ENCODING = "json"

type Node interface {
	// Type returns a static name of the type of the node for debugging purposes.
	Type() string
}

// DefaultNode is a node that does nothing, but provides a base for other nodes to compose on top of.
type DefaultNode struct{}

func (d *DefaultNode) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error { return nil }
func (d *DefaultNode) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	return nil
}

// NodeConstructor is a function that takes a raw graph node and turns it into a node that can actually process things.
type NodeConstructor = func(n node) (Node, error)

var nodeRegistry = map[string]NodeConstructor{
	"":       func(n node) (Node, error) { return &AnchorNode{}, nil },
	"stdout": NewFileNotifierNode,
	"stderr": NewFileNotifierNode,
	"file":   NewFileNotifierNode,
}

// LookupNode takes a node type name and returns a constructor that can be used to make nodes of that name.
func LookupNode(name string) NodeConstructor {
	return nodeRegistry[name]
}

// AnchorNode is the default node type, if nothing else is specified. They do nothing except
// act as anchor points for Links to allow splitting one or more incoming links into one or more outgoing ones.
type AnchorNode struct {
	*DefaultNode
}

func (a *AnchorNode) Type() string {
	return "anchor"
}

// FileNotifierNode represents a node that can output alerts to a Writer.
type FileNotifierNode struct {
	*DefaultNode
	encoder encoding.Encoder
	file    io.WriteCloser
}

func NewFileNotifierNode(n node) (Node, error) {
	encodingName := DEFAULT_ENCODING
	if enc, ok := n.attrs["encoding"]; ok {
		encodingName = enc
	}

	encoder := encoding.LookupEncoding(encodingName)
	if encoder == nil {
		return nil, fmt.Errorf("invalid encoding: %q", encodingName)
	}

	switch n.attrs["type"] {
	case "stdout":
		return &FileNotifierNode{
			encoder: encoder,
			file:    os.Stdout,
		}, nil
	case "stderr":
		return &FileNotifierNode{
			encoder: encoder,
			file:    os.Stderr,
		}, nil
	case "", "file":
		fileName := n.attrs["path"]
		if fileName == "" {
			return nil, errors.New("missing `path` in file node")
		}

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %q in file node: %w", fileName, err)
		}

		return &FileNotifierNode{
			encoder: encoder,
			file:    file,
		}, nil
	default:
		return nil, fmt.Errorf("invalid type for file node: %q", n.attrs["type"])
	}
}

func (f *FileNotifierNode) Type() string {
	return "file"
}

func (f *FileNotifierNode) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	var lastError error
	for _, alert := range alerts {
		bytes, err := f.encoder.Marshal(alert)
		if err != nil {
			lastError = multierror.Append(lastError, err)
			continue
		}

		if _, err := f.file.Write(bytes); err != nil {
			lastError = multierror.Append(lastError, err)
		}
	}

	return lastError
}

func (f *FileNotifierNode) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	var lastError error
	for _, silence := range silences {
		bytes, err := f.encoder.Marshal(silence)
		if err != nil {
			lastError = multierror.Append(lastError, err)
			continue
		}

		if _, err := f.file.Write(bytes); err != nil {
			lastError = multierror.Append(lastError, err)
		}
	}

	return lastError
}
