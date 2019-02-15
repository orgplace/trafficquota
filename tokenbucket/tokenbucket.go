// Package tokenbucket provides an interface of token bucket and its implementation.
package tokenbucket

import "time"

//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -write_package_comment=false -destination=mock_${GOFILE}

const (
	// DefaultInterval is default interval to fill bucket.
	DefaultInterval = 10 * time.Millisecond
	// DefaultRate is filled tokens per seccond.
	DefaultRate int32 = 100
	// DefaultBucketSize is default size of bucket.
	DefaultBucketSize = DefaultRate / 5
)

// Config is a configuration of TokenBucket.
type Config interface {
	// Rate returns a number of filled tokens.
	Rate(partitionKey, clusteringKey string) int32
	// Overflow returns true when .
	Overflow(partitionKey, clusteringKey string, tokens int32) bool
}

// TokenBucket is an algorithm used to control network traffic.
// This interface provices goroutine-safe methods.
type TokenBucket interface {
	Fill()
	Take(partitionKey string, clusteringKeys []string) (bool, error)
}
