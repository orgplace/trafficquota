package config

import (
	"time"

	"github.com/orgplace/trafficquota/tokenbucket"
)

// FileContent is a content of config file.
type FileContent struct {
	TokenBucket tokenBucketConfig
}

type tokenBucketConfig struct {
	Interval duration

	Default tokenbucket.BucketOption
	Chunks  map[string]*tokenbucket.ChunkOption
}

func (c *tokenBucketConfig) AsOption() *tokenbucket.Option {
	return &tokenbucket.Option{
		Interval: c.Interval.Duration,
		Default:  c.Default,
		Chunks:   c.Chunks,
	}
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
