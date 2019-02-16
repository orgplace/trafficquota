package tokenbucket

import "time"

type defaultConfig struct{}

var (
	// DefaultConfig is the default configuration.
	DefaultConfig = defaultConfig{}

	defaultRatePerInterval = toFilled(DefaultRate, DefaultInterval)
)

func (c defaultConfig) Rate(_, _ string) int32 {
	return defaultRatePerInterval
}

func (c defaultConfig) Overflow(_, _ string, tokens int32) bool {
	return DefaultBucketSize < tokens
}

func toFilled(filledPerSec int32, interval time.Duration) int32 {
	return filledPerSec / int32(time.Second/interval)
}
