package unmarshal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalOpts is a struct of the options that can be passed to UnmarshalConfig.
type UnmarshalOpts struct {
	DisallowUnknownFields bool
}

// UnmarshalConfig unmarshals a struct from a map of strings to strings using the tags on the exported fields of the struct.
// When unmarshaling MaybeFiles, it checks both the field name, and the field name suffixed with _file for literal values and files specifically.
func UnmarshalConfig(data map[string]string, v interface{}, opts UnmarshalOpts) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("UnmarshalConfig: invalid argument, must be a non-nil pointer")
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("UnmarshalConfig: invalid argument, must be a struct pointer")
	}

	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// Skip unexported fields.
		if !fieldValue.CanSet() {
			continue
		}

		fieldName := field.Name
		fieldType := field.Type
		fieldValueTag := fieldValue.Addr().Interface()

		if tag := field.Tag.Get("config"); tag != "" {
			fieldName = tag
		}

		if fieldName == "-" {
			continue
		}

		// For MaybeFile and MaybeSecretFile s, check if the field name suffixed with _file is specified.
		if fieldType == reflect.TypeOf(MaybeFile{}) || fieldType == reflect.TypeOf(MaybeSecretFile{}) || fieldType == reflect.TypeOf(&MaybeFile{}) || fieldType == reflect.TypeOf(&MaybeSecretFile{}) {
			fileFieldName := fmt.Sprintf("%s_file", fieldName)

			if _, ok := data[fileFieldName]; ok {
				if _, ok := data[fieldName]; ok {
					return fmt.Errorf("UnmarshalConfig: field %s cannot be specified as both a literal value and a file", fieldName)
				}
			}

			if fieldType == reflect.TypeOf(MaybeFile{}) || fieldType == reflect.TypeOf(&MaybeFile{}) {
				val, err := NewMaybeFile(data[fileFieldName], data[fieldName])
				if err != nil {
					return err
				}

				if fieldType == reflect.TypeOf(MaybeFile{}) {
					fieldValue.Set(reflect.ValueOf(*val))
				} else {
					fieldValue.Set(reflect.ValueOf(val))
				}
			} else if fieldType == reflect.TypeOf(MaybeSecretFile{}) || fieldType == reflect.TypeOf(&MaybeSecretFile{}) {
				val, err := NewMaybeSecretFile(data[fileFieldName], Secret(data[fieldName]))
				if err != nil {
					return err
				}

				if fieldType == reflect.TypeOf(MaybeSecretFile{}) {
					fieldValue.Set(reflect.ValueOf(*val))
				} else {
					fieldValue.Set(reflect.ValueOf(val))
				}
			}

			continue
		}

		fieldValueStr, ok := data[fieldName]
		if !ok {
			if _, ok := field.Tag.Lookup("required"); ok {
				return fmt.Errorf("UnmarshalConfig: field %s is required but not found in the config", field.Name)
			}
			continue
		}

		if err := unmarshalValue(fieldValueStr, fieldValueTag); err != nil {
			return err
		}

		delete(data, fieldName)
	}

	if opts.DisallowUnknownFields && len(data) > 0 {
		return fmt.Errorf("found extra fields while unmarshaling: %v", data)
	}

	return nil
}

func unmarshalValue(valueStr string, v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got %T", v)
	}

	switch v := v.(type) {
	case *string:
		*v = valueStr
	case **string:
		*v = &valueStr
	case *bool, **bool:
		boolValue, err := strconv.ParseBool(valueStr)
		if err != nil {
			return err
		}

		if v, ok := v.(*bool); ok {
			*v = boolValue
		}

		if v, ok := v.(**bool); ok {
			*v = &boolValue
		}
	case *int, **int:
		intValue, err := strconv.Atoi(valueStr)
		if err != nil {
			return err
		}

		if v, ok := v.(*int); ok {
			*v = intValue
		}

		if v, ok := v.(**int); ok {
			*v = &intValue
		}
	case *float64, **float64:
		floatValue, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return err
		}

		if v, ok := v.(*float64); ok {
			*v = floatValue
		}

		if v, ok := v.(**float64); ok {
			*v = &floatValue
		}
	case *float32, **float32:
		floatValue, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			return err
		}

		if v, ok := v.(*float64); ok {
			*v = floatValue
		}

		if v, ok := v.(**float64); ok {
			*v = &floatValue
		}
	case *[]string:
		strSlice := strings.Split(valueStr, ",")
		*v = strSlice
	default:
		return fmt.Errorf("UnmarshalConfig: unsupported field type: %T", v)
	}

	return nil
}
