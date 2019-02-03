package tokenbucket

import "time"

const (
	DefaultInterval   = 10 * time.Millisecond
	DefaultRate       = 100
	DefaultBucketSize = 20
)

type TokenBucket interface {
	Take(partitionKey string, clusteringKeys []string) (bool, error)
}
