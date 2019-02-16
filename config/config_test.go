package config

import (
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap/zapcore"
)

func Test_getLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  zapcore.Level
	}{
		{
			name:  "empty",
			key:   "EMPTY_KEY",
			value: "",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "info",
			key:   "INFO_KEY",
			value: "INFO",
			want:  zapcore.InfoLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.key, tt.value)
			if got := getLogLevel(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configGetter.getLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLogLevel_invalid(t *testing.T) {
	os.Setenv("LOG_LEVEL", "invalid log level")

	var err interface{}
	func() {
		defer func() {
			err = recover()
		}()
		getLogLevel("LOG_LEVEL")
	}()

	if err == nil {
		t.Errorf("configGetter.getLogLevel() recovered = %v", err)
	}
}

func Test_getEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		want         string
	}{
		{
			name:         "default",
			key:          "DEFAULT",
			value:        "",
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "value",
			key:          "VALUE",
			value:        "val",
			defaultValue: "default",
			want:         "val",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.key, tt.value)
			if got := getEnv(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("configGetter.getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
