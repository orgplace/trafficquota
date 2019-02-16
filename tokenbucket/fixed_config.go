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
	defaultSize   int32
	bucketSize    map[string]int32
	minBucketSize int32
}

type fixedChunkRateConfig struct {
	defultRate int32
	bucketRate map[string]int32
}

func (o *BucketOption) getSize(def int32) int32 {
	if o.Banned {
		return 0
	}
	if o.Size == 0 {
		return def
	}
	return o.Size
}

// NewFixedConfig constructs a new configuration.
// The fixed config is read only after creation.
func NewFixedConfig(o *Option) Config {
	interval := o.Interval
	if interval == 0 {
		interval = DefaultInterval
	}

	defaultSize := o.Default.getSize(DefaultBucketSize)

	n := len(o.Chunks)
	chunkSize := make(map[string]*fixedChunkSizeConfig, n)
	minChunkSize := defaultSize
	chunkRate := make(map[string]*fixedChunkRateConfig, n)
	for k, v := range o.Chunks {
		cs, min := buildFixedChunkSizeConfig(v, v.Default.getSize(defaultSize))
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

func buildFixedChunkSizeConfig(o *ChunkOption, defaultSize int32) (*fixedChunkSizeConfig, int32) {
	bucketSize := make(map[string]int32, len(o.Chunk))
	minSize := defaultSize
	for k, v := range o.Chunk {
		s := v.getSize(defaultSize)
		bucketSize[k] = s
		minSize = minInt32(minSize, s)
	}

	return &fixedChunkSizeConfig{
		defaultSize:   defaultSize,
		bucketSize:    bucketSize,
		minBucketSize: minSize,
	}, minSize
}

func buildFixedChunkRateConfig(o *ChunkOption, interval time.Duration) *fixedChunkRateConfig {
	bucketRate := make(map[string]int32, len(o.Chunk))
	for k, v := range o.Chunk {
		bucketRate[k] = toFilled(v.Rate, interval)
	}

	return &fixedChunkRateConfig{
		defultRate: toFilled(o.Default.Rate, interval),
		bucketRate: bucketRate,
	}
}

func (c *fixedConfig) Overflow(chunkKey, bucketKey string, tokens int32) bool {
	if tokens <= c.minChunkSize {
		return false
	}

	chunkSize, ok := c.chunkSize[chunkKey]
	if !ok {
		return c.defaultSize < tokens
	}

	if tokens <= chunkSize.minBucketSize {
		return false
	}

	bucketSize, ok := chunkSize.bucketSize[bucketKey]
	if !ok {
		return chunkSize.defaultSize < tokens
	}
	return bucketSize < tokens
}

func (c *fixedConfig) Rate(chunkKey, bucketKey string) int32 {
	chunkRate, ok := c.chunkRate[chunkKey]
	if !ok {
		return c.defaultRate
	}

	rate, ok := chunkRate.bucketRate[bucketKey]
	if !ok {
		return chunkRate.defultRate
	}
	return rate
}
