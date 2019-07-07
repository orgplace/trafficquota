// Package tokenbucket provides an interface of token bucket and its implementation.
package tokenbucket

import "time"

//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -write_package_comment=false -destination=mock_${GOFILE}

const (
	// DefaultTimeSlice is default interval to fill bucket.
	DefaultTimeSlice = 10 * time.Millisecond
	// DefaultGCInterval is default interval to collect garbege token.
	DefaultGCInterval = 1 * time.Second
	// DefaultRate is filled tokens per seccond.
	DefaultRate int32 = 100
	// DefaultBucketSize is default size of bucket.
	DefaultBucketSize = DefaultRate / 5
)

// Config is a configuration of TokenBucket.
type Config interface {
	// Rate returns a number of filled tokens.
	Rate(chunkKey, bucketKey string) int32
	// Overflow returns true when .
	Overflow(chunkKey, bucketKey string, tokens int32) bool
}

// TokenBucket is an algorithm used to control network traffic.
// This interface provices goroutine-safe methods.
type TokenBucket interface {
	Take(chunkKey string, bucketKeys []string) (bool, error)
}

// TimeSliceTokenBucket is an TokenBucket based on time slices.
type TimeSliceTokenBucket interface {
	TokenBucket
	Fill()
}

// TimestampTokenBucket is an TokenBucket based on timestamp.
type TimestampTokenBucket interface {
	TokenBucket
	GC()
}
