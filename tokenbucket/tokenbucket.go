package tokenbucket

import "time"

//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=mock_${GOFILE}

const (
	// DefaultInterval is default interval to fill bucket.
	DefaultInterval = 10 * time.Millisecond
	// DefaultRate is filled tokens per seccond.
	DefaultRate = 100
	// DefaultBucketSize is default size of bucket.
	DefaultBucketSize = DefaultRate / 5
)

// TokenBucket is an algorithm used to control network traffic.
// This interface provices goroutine-safe methods.
type TokenBucket interface {
	Fill()
	Take(partitionKey string, clusteringKeys []string) (bool, error)
}
