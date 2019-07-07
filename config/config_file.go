package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/orgplace/trafficquota/tokenbucket"
)

// FileContent is a content of config file.
type FileContent struct {
	Strategy    Strategy
	TokenBucket tokenBucketConfig
}

type fileContent struct {
	Strategy    strategy
	TokenBucket tokenBucketConfig
}

type strategy Strategy

func (s *strategy) UnmarshalText(text []byte) error {
	err := (*Strategy)(s).UnmarshalText(text)
	if err != nil {
		*s = strategy(StrategyTimeSlice)
	}
	return nil
}

type tokenBucketConfig struct {
	Interval duration

	Default tokenbucket.BucketOption
	Chunks  map[string]*tokenbucket.ChunkOption
}

type duration time.Duration

func (d *duration) Duration(def time.Duration) time.Duration {
	if *d == duration(0) {
		return def
	}
	return time.Duration(*d)
}

func (d *duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))
	*d = duration(dur)
	return err
}

func LoadFile(path string) (*FileContent, error) {
	var result fileContent
	switch _, err := toml.DecodeFile(path, &result); err.(type) {
	case *os.PathError, nil:
		// nothing to do
	default:
		return nil, err
	}

	return &FileContent{
		Strategy:    Strategy(result.Strategy),
		TokenBucket: result.TokenBucket,
	}, nil
}

func (c *tokenBucketConfig) AsOption(defaultInterval time.Duration) *tokenbucket.Option {
	return &tokenbucket.Option{
		Interval: c.Interval.Duration(defaultInterval),
		Default:  c.Default,
		Chunks:   c.Chunks,
	}
}
