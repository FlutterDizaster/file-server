package configloader

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func parseFile(val reflect.Value, filePath string) error {
	if filePath == "" {
		return nil
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("could not open config file: " + err.Error())
	}
	defer file.Close()

	// Decode file to map
	var data map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return errors.New("could not decode config file: " + err.Error())
	}

	// Populate struct
	return populateStruct(val, data)
}

func populateStruct(val reflect.Value, data map[string]interface{}) error {
	typ := val.Type()
	var combinedError error

	for i := range typ.NumField() {
		field := typ.Field(i)
		value := val.Field(i)

		// Obtain field name
		name := field.Tag.Get("name")
		if name == "" {
			name = strings.ToLower(field.Name)
		}

		// Set field
		if fieldValue, exists := data[name]; exists {
			strValue, ok := fieldValue.(string)
			if !ok {
				strValue = ""
			}
			if err := setField(value, strValue); err != nil {
				fieldError := errors.New("could not set field " + name + ": " + err.Error())
				combinedError = errors.Join(combinedError, fieldError)
			}
		}
	}
	return combinedError
}

func setField(field reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	//nolint:exhaustive // all supported types
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, field.Type().Bits())
		if err != nil {
			return errors.New("expected integer value")
		}
		field.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, field.Type().Bits())
		if err != nil {
			return errors.New("expected unsigned integer value")
		}
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return errors.New("expected float value")
		}
		field.SetFloat(floatVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return errors.New("expected boolean value")
		}
		field.SetBool(boolVal)

	case reflect.Struct:
		var mapVal map[string]interface{}
		if err := json.Unmarshal([]byte(value), &mapVal); err != nil {
			return errors.New("expected JSON object for struct field")
		}
		return populateStruct(field, mapVal)

	default:
		return errors.New("unsupported field type: " + field.Kind().String())
	}
	return nil
}
