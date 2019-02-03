package config

import (
	"go.uber.org/zap/zapcore"
)

var (
	DevelopMode bool
	LogLevel    zapcore.Level

	Listen string
)

func init() {
	DevelopMode = true
	LogLevel = zapcore.DebugLevel

	Listen = "unix:/tmp/test.sock"
	// Listen = "localhost:3895"
}
