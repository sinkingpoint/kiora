package unmarshal_test

import (
	"testing"

	"github.com/sinkingpoint/kiora/lib/kiora/config/unmarshal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalConfig(t *testing.T) {
	data := map[string]string{
		"field1": "value1",
		"field2": "42",
		"field3": "3.14",
		"field4": "true",
		"field5": "value1,value2,value3",
		"field6": "foo",
	}

	type Config struct {
		Field1 string                    `config:"field1"`
		Field2 int                       `config:"field2"`
		Field3 float64                   `config:"field3"`
		Field4 bool                      `config:"field4"`
		Field5 []string                  `config:"field5"`
		Field6 unmarshal.MaybeSecretFile `config:"field6"`
		Field7 *unmarshal.MaybeFile      `config:"field6"`
	}

	file, err := unmarshal.NewMaybeFile("", "foo")
	require.NoError(t, err)

	secretFile, err := unmarshal.NewMaybeSecretFile("", "foo")
	require.NoError(t, err)

	expected := Config{
		Field1: "value1",
		Field2: 42,
		Field3: 3.14,
		Field4: true,
		Field5: []string{"value1", "value2", "value3"},
		Field6: *secretFile,
		Field7: file,
	}

	var config Config
	require.NoError(t, unmarshal.UnmarshalConfig(data, &config, unmarshal.UnmarshalOpts{}), "Unexpected error")

	assert.Equal(t, expected, config, "Incorrect unmarshalled config")
}

func TestUnmarshalConfig_MissingRequiredField(t *testing.T) {
	data := map[string]string{
		"field1": "value1",
		"field2": "42",
	}

	type Config struct {
		Field1 string `config:"field1" required:"true"`
		Field2 int    `config:"field2" required:"true"`
		Field3 string `config:"field3" required:"true"`
	}

	var config Config
	err := unmarshal.UnmarshalConfig(data, &config, unmarshal.UnmarshalOpts{})
	assert.Error(t, err, "Expected error")
	expectedError := "UnmarshalConfig: field Field3 is required but not found in the config"
	assert.EqualError(t, err, expectedError, "Incorrect error message")
}

func TestUnmarshalConfig_DisallowUnknownFields(t *testing.T) {
	data := map[string]string{
		"field1":     "value1",
		"field2":     "42",
		"unexpected": "true",
	}

	type Config struct {
		Field1 string `config:"field1" required:"true"`
		Field2 int    `config:"field2" required:"true"`
	}

	var config Config
	err := unmarshal.UnmarshalConfig(data, &config, unmarshal.UnmarshalOpts{
		DisallowUnknownFields: true,
	})
	assert.Error(t, err, "Expected error")
}

func TestUnmarshalConfig_DisallowBothFileAndLiteral(t *testing.T) {
	data := map[string]string{
		"field1":      "value1",
		"field1_file": "./foo.txt",
	}

	type Config struct {
		Field1 *unmarshal.MaybeFile `config:"field1" required:"true"`
	}

	var config Config
	err := unmarshal.UnmarshalConfig(data, &config, unmarshal.UnmarshalOpts{})
	assert.Error(t, err, "Expected error")
}
