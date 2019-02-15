package tokenbucket

import "time"

type defaultConfig struct{}

// DefaultConfig is the default configuration.
var DefaultConfig = defaultConfig{}

func (c defaultConfig) Rate(_, _ string) int32 {
	return DefaultRate / int32(time.Second/DefaultInterval)
}

func (c defaultConfig) Overflow(_, _ string, tokens int32) bool {
	return DefaultBucketSize < tokens
}
