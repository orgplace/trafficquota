package config

import (
	"reflect"
	"testing"

	"go.uber.org/zap/zapcore"
)

func mockConfigGetter(t *testing.T, expectedKey, value string) configGetter {
	return configGetter(func(key string) string {
		if key != expectedKey {
			t.Errorf("got key %v, want %v", key, expectedKey)
		}
		return value
	})
}

func Test_configGetter_getLogLevel(t *testing.T) {
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
			g := mockConfigGetter(t, tt.key, tt.value)
			if got := g.getLogLevel(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configGetter.getLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGetter_getEnv(t *testing.T) {
	tests := []struct {
		name         string
		g            configGetter
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
			g := mockConfigGetter(t, tt.key, tt.value)
			if got := g.getEnv(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("configGetter.getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
