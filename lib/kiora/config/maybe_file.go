package config

import (
	"os"

	"github.com/pkg/errors"
)

// MaybeFile is a value that can be loaded from an external file, or specified directly.
type MaybeFile struct {
	// path is the path to the file.
	path string

	// value is the contents of the file, or the literal value.
	value string
}

// NewMaybeFile creates a new MaybeFile, loading the value from the file if path is not empty.
func NewMaybeFile(path, value string) (*MaybeFile, error) {
	if path != "" && value != "" {
		return nil, errors.New("cannot specify both path and value")
	}

	if path != "" {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		value = string(bytes)
	}

	return &MaybeFile{
		path:  path,
		value: value,
	}, nil
}

func (m *MaybeFile) Value() string {
	return m.value
}

// MaybeSecretFile is a MaybeFile that redacts the value.
type MaybeSecretFile struct {
	// path is the path to the file.
	path string

	// value is the contents of the file, or the literal value.
	value Secret
}

func NewMaybeSecretFile(path string, value Secret) (*MaybeSecretFile, error) {
	if path != "" && value != "" {
		return nil, errors.New("cannot specify both path and value")
	}

	if path != "" {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		value = Secret(string(bytes))
	}

	return &MaybeSecretFile{
		path:  path,
		value: value,
	}, nil
}

func (m *MaybeSecretFile) Value() Secret {
	return m.value
}
