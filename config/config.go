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
	g := configGetter(os.Getenv)
	LogLevel = g.getLogLevel("LOG_LEVEL")

	Listen = g.getEnv("LISTEN", "127.0.0.1:3895")
}

type configGetter func(string) string

func (g configGetter) getLogLevel(key string) zapcore.Level {
	val := g(key)
	if val == "" {
		return zapcore.DebugLevel
	}

	var l zapcore.Level
	if err := l.Set(val); err != nil {
		panic(err)
	}
	return l
}

func (g configGetter) getEnv(key, def string) string {
	val := g(key)
	if val == "" {
		return def
	}
	return val
}
