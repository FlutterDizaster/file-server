package configloader

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Host   string `name:"host"`
	Port   int    `name:"port"`
	Debug  bool   `name:"debug"`
	Nested struct {
		Enabled bool `name:"enabled"`
		Level   int  `name:"level"`
	} `name:"nested"`
}

func Test_parseFile(t *testing.T) {
	type test struct {
		name      string
		fileData  string
		filePath  string
		expected  Config
		expectErr bool
	}

	tests := []test{
		{
			name: "Valid config file",
			//nolint:lll // test data
			fileData: `{"host": "localhost", "port": "8080", "debug": "true", "nested": "{\"enabled\": \"true\", \"level\": \"2\"}"}`,
			filePath: "config_test.json",
			expected: Config{Host: "localhost", Port: 8080, Debug: true, Nested: struct {
				Enabled bool `name:"enabled"`
				Level   int  `name:"level"`
			}{Enabled: true, Level: 2}},
			expectErr: false,
		},
		{
			name:      "Invalid JSON format",
			fileData:  `{"host": "localhost", "port": "invalid", "debug": "true"`,
			filePath:  "config_test.json",
			expected:  Config{},
			expectErr: true,
		},
		{
			name:      "Empty file path",
			fileData:  "",
			filePath:  "",
			expected:  Config{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.filePath
			if filePath != "" {
				tmpFile, err := os.CreateTemp("", filePath)
				require.NoError(t, err)
				defer os.Remove(tmpFile.Name())

				if tt.fileData != "" {
					_, ferr := tmpFile.WriteString(tt.fileData)
					require.NoError(t, ferr)
				}

				require.NoError(t, tmpFile.Close())

				filePath = tmpFile.Name()
			}

			var cfg Config
			val := reflect.ValueOf(&cfg).Elem()

			err := parseFile(val, filePath)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, cfg)
			}
		})
	}
}

func Test_populateStruct(t *testing.T) {
	type test struct {
		name      string
		data      map[string]interface{}
		expected  Config
		expectErr bool
	}

	tests := []test{
		{
			name: "Valid data",
			data: map[string]interface{}{
				"host":   "localhost",
				"port":   "8080",
				"debug":  "true",
				"nested": `{"enabled": "true", "level": "2"}`,
			},
			expected: Config{
				Host:  "localhost",
				Port:  8080,
				Debug: true,
				Nested: struct {
					Enabled bool `name:"enabled"`
					Level   int  `name:"level"`
				}{Enabled: true, Level: 2},
			},
			expectErr: false,
		},
		{
			name:      "Type mismatch",
			data:      map[string]interface{}{"port": "not_a_number"},
			expected:  Config{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			val := reflect.ValueOf(&cfg).Elem()

			err := populateStruct(val, tt.data)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, cfg)
			}
		})
	}
}

func Test_setField(t *testing.T) {
	type test struct {
		name      string
		value     string
		fieldKind reflect.Kind
		expectErr bool
		expected  interface{}
	}

	tests := []test{
		{
			name:      "Set string field",
			value:     "localhost",
			fieldKind: reflect.String,
			expected:  "localhost",
		},
		{
			name:      "Set int field",
			value:     "8080",
			fieldKind: reflect.Int,
			expected:  8080,
		},
		{
			name:      "Set bool field",
			value:     "true",
			fieldKind: reflect.Bool,
			expected:  true,
		},
		{
			name:      "Invalid int value",
			value:     "invalid",
			fieldKind: reflect.Int,
			expectErr: true,
			expected:  int64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.New(reflect.TypeOf(tt.expected)).Elem()
			field.Set(reflect.Zero(reflect.TypeOf(tt.expected)))

			err := setField(field, tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, field.Interface())
			}
		})
	}
}
