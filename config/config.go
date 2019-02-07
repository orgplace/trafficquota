package config

import (
	"os"

	"go.uber.org/zap/zapcore"
)

var (
	LogLevel zapcore.Level

	Listen string
)

func init() {
	LogLevel = getLogLevel()

	// Listen = "unix:/tmp/test.sock"
	Listen = getEnv("LISTEN", "127.0.0.1:3895")
}

func getLogLevel() zapcore.Level {
	val := os.Getenv("LOG_LEVEL")
	if val == "" {
		return zapcore.DebugLevel
	}

	var l zapcore.Level
	if err := l.Set(val); err != nil {
		panic(err)
	}
	return l
}

func getEnv(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
