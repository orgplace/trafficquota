package config

import (
	"os"

	"go.uber.org/zap/zapcore"
)

var (
	// LogLevel is a log level.
	LogLevel zapcore.Level

	// Listen is address to listen.
	// It is must be host:port or unix:/path/of/sock .
	Listen string
)

func init() {
	LogLevel = getLogLevel("LOG_LEVEL")

	Listen = getEnv("LISTEN", "0.0.0.0:3895")
}

func getLogLevel(key string) zapcore.Level {
	val := os.Getenv(key)
	if val == "" {
		return zapcore.DebugLevel
	}

	var l zapcore.Level
	if err := l.Set(val); err != nil {
		panic(err)
	}
	return l
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
