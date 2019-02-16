package tokenbucket

import "time"

// Option is an option of tokenbucket
type Option struct {
	Interval time.Duration

	Default BucketOption
	Chunks  map[string]*ChunkOption
}

func (o *Option) getInterval() time.Duration {
	if o.Interval == 0 {
		return DefaultInterval
	}
	return o.Interval
}

// ChunkOption is an option of chunk
type ChunkOption struct {
	Default BucketOption
	Chunk   map[string]*BucketOption
}

// BucketOption is an option of bucket
type BucketOption struct {
	// Banned means whether the bucket is banned.
	// When the value is true,
	// size will be 0 and any requets for the bucket will be rejected.
	// Otherwise, the default size is used.
	Banned bool
	// Size is size of the bucket.
	Size int32
	// Rate is rate of the bucket per second.
	Rate int32
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

func (o *BucketOption) getRate(def int32) int32 {
	if o.Rate == 0 {
		return def
	}
	return o.Rate
}
