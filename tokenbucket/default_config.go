package tokenbucket

import "time"

type defaultConfig struct{}

var (
	// DefaultConfig is the default configuration.
	DefaultConfig = defaultConfig{}

	defaultRatePerInterval = toFilled(DefaultRate, DefaultTimeSlice)
)

func (c defaultConfig) Rate(_, _ string) int32 {
	return defaultRatePerInterval
}

func (c defaultConfig) Overflow(_, _ string, tokens int32) bool {
	return DefaultBucketSize < tokens
}

func toFilled(filledPerSec int32, interval time.Duration) int32 {
	return int32(int64(filledPerSec) * int64(interval) / int64(time.Second))
}
