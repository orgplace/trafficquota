package tokenbucket

import "time"

type fixedConfig struct {
	defaultSize int32
	chunkSize   map[string]*fixedChunkSizeConfig
	// minChunkSize is used to cut size search.
	minChunkSize int32

	defaultRate int32
	chunkRate   map[string]*fixedChunkRateConfig
}

type fixedChunkSizeConfig struct {
	defultSize    int32
	bucketSize    map[string]int32
	minBucketSize int32
}

type fixedChunkRateConfig struct {
	defultRate int32
	bucketRate map[string]int32
}

// NewFixedConfig constructs a new configuration.
// The fixed config is read only after creation.
func NewFixedConfig(o *Option) Config {
	interval := o.Interval
	if interval == 0 {
		interval = DefaultInterval
	}

	defaultSize := int32(0)
	if !o.Default.Banned {
		defaultSize = o.Default.Size
	}

	n := len(o.Partitions)
	chunkSize := make(map[string]*fixedChunkSizeConfig, n)
	minChunkSize := defaultSize
	chunkRate := make(map[string]*fixedChunkRateConfig, n)
	for k, v := range o.Partitions {
		cs, min := buildFixedChunkSizeConfig(v, o.Default.Banned)
		chunkSize[k] = cs
		minChunkSize = minInt32(minChunkSize, min)
		chunkRate[k] = buildFixedChunkRateConfig(v, interval)
	}

	return &fixedConfig{
		defaultSize:  defaultSize,
		chunkSize:    chunkSize,
		minChunkSize: minChunkSize,

		defaultRate: toFilled(o.Default.Rate, interval),
		chunkRate:   chunkRate,
	}
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func buildFixedChunkSizeConfig(o *ChunkOption, bannedDefault bool) (*fixedChunkSizeConfig, int32) {
	defultSize := int32(0)
	if !bannedDefault {
		defultSize = o.Default.Size
	}

	bucketSize := make(map[string]int32, len(o.Buckets))
	minSize := defultSize
	for k, v := range o.Buckets {
		s := int32(0)
		if !v.Banned {
			s = v.Size
		}
		bucketSize[k] = s
		minSize = minInt32(minSize, s)
	}

	return &fixedChunkSizeConfig{
		defultSize:    defultSize,
		bucketSize:    bucketSize,
		minBucketSize: minSize,
	}, minSize
}

func buildFixedChunkRateConfig(o *ChunkOption, interval time.Duration) *fixedChunkRateConfig {
	bucketRate := make(map[string]int32, len(o.Buckets))
	for k, v := range o.Buckets {
		bucketRate[k] = toFilled(v.Rate, interval)
	}

	return &fixedChunkRateConfig{
		defultRate: toFilled(o.Default.Rate, interval),
		bucketRate: bucketRate,
	}
}

func (c *fixedConfig) Overflow(partitionKey, chunkKey string, tokens int32) bool {
	if tokens <= c.minChunkSize {
		return false
	}

	chunkSize, ok := c.chunkSize[partitionKey]
	if !ok {
		return c.defaultSize < tokens
	}

	if tokens <= chunkSize.minBucketSize {
		return false
	}

	bucketSize, ok := chunkSize.bucketSize[chunkKey]
	if !ok {
		return chunkSize.defultSize < tokens
	}
	return bucketSize < tokens
}

func (c *fixedConfig) Rate(partitionKey, chunkKey string) int32 {
	chunkRate, ok := c.chunkRate[partitionKey]
	if !ok {
		return c.defaultRate
	}

	rate, ok := chunkRate.bucketRate[chunkKey]
	if !ok {
		return chunkRate.defultRate
	}
	return rate
}
