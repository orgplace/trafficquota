package tokenbucket

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type inMemoryTokenBucket struct {
	shards bucketShards
}

// NewInMemoryTokenBucket constructs a in-memory TokenBucket
func NewInMemoryTokenBucket() TokenBucket {
	b := &inMemoryTokenBucket{
		shards: *newBucketShard(),
	}

	go func() {
		c := time.Tick(DefaultInterval)
		for range c {
			b.shards.each(func(partitionKey string, b *buckets) {
				b.fill()
				// TODO: Delete empty bucket
			})
		}
	}()

	return b
}

func (b *inMemoryTokenBucket) Take(partitionKey string, clusteringKeys []string) (bool, error) {
	buckets := b.shards.loadOrStore(partitionKey, newBuckets())
	for _, clusteringKey := range clusteringKeys {
		if ok := buckets.take(partitionKey, clusteringKey); !ok {
			return false, nil
		}
	}
	return true, nil
}

type bucketShards sync.Map

func newBucketShard() *bucketShards {
	return &bucketShards{}
}

func (s *bucketShards) each(f func(key string, value *buckets)) {
	(*sync.Map)(s).Range(func(key, value interface{}) bool {
		f(key.(string), value.(*buckets))
		return false
	})
}

func (s *bucketShards) loadOrStore(key string, b *buckets) *buckets {
	v, _ := (*sync.Map)(s).LoadOrStore(key, b)
	return v.(*buckets)
}

type buckets struct {
	// Map of [key] => [took token].
	// Took token is start from 0.
	// Then, ([took token] + 1) access was permitted.
	m *sync.Map
}

func newBuckets() *buckets {
	return &buckets{
		m: &sync.Map{},
	}
}

// expungedBucket means that bucket was expunged.
// After swap to this value, the bucket must be deleted from buckets.
const expungedBucket = int32(math.MinInt32)

func (b *buckets) fill() {
	n := int32(DefaultRate / int32(time.Second/DefaultInterval))
	b.m.Range(func(key, value interface{}) bool {
		p := value.(*int32)

		for {
			current := atomic.LoadInt32(p)
			if current == expungedBucket {
				break
			}

			next := current - n
			if next < 0 {
				if atomic.CompareAndSwapInt32(p, current, expungedBucket) {
					b.m.Delete(key)
					break
				}
			} else if atomic.CompareAndSwapInt32(p, current, next) {
				break
			}
		}

		return false
	})
}

func (b *buckets) take(partitionKey, clusteringKey string) bool {
LOAD_OR_NEW_LOOP:
	for {
		value, loaded := b.loadOrNew(clusteringKey)
		if !loaded {
			return true
		}

		for {
			current := atomic.LoadInt32(value)
			if current == expungedBucket {
				continue LOAD_OR_NEW_LOOP
			}

			next := current + 1
			if DefaultBucketSize <= next {
				// TODO: load configured size
				//if configuredBucketSize <= current {
				return false
				//}
			}

			if atomic.CompareAndSwapInt32(value, current, maxInt32(next, 0)) {
				return true
			}
		}
	}
}

func (b *buckets) loadOrNew(clusteringKey string) (*int32, bool) {
	newValue := new(int32)
	for {
		p, loaded := b.m.LoadOrStore(clusteringKey, newValue)
		if !loaded {
			return newValue, false
		}

		value := p.(*int32)
		current := atomic.LoadInt32(value)
		if current != expungedBucket {
			return value, true
		}
	}
}

func maxInt32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}
