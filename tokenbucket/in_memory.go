package tokenbucket

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type inMemoryTokenBucket struct {
	shards sync.Map
}

// NewInMemoryTokenBucket constructs a in-memory TokenBucket
func NewInMemoryTokenBucket() TokenBucket {
	tb := &inMemoryTokenBucket{
		shards: sync.Map{},
	}

	go func() {
		c := time.Tick(DefaultInterval)
		for range c {
			tb.shards.Range(func(key, value interface{}) bool {
				b := value.(*buckets)

				if b.fill() {
					b.mu.Lock()
					if b.empty() {
						b.expunged = true
					}
					b.mu.Unlock()
					tb.shards.Delete(key)
				}

				return true
			})
		}
	}()

	return tb
}

func (tb *inMemoryTokenBucket) Take(partitionKey string, clusteringKeys []string) (bool, error) {
	newValue := newBuckets()
	for {
		value, _ := tb.shards.LoadOrStore(partitionKey, newValue)
		b := value.(*buckets)

		b.mu.Lock()
		if b.expunged {
			b.mu.Unlock()
			continue
		}

		for _, clusteringKey := range clusteringKeys {
			if ok := b.take(partitionKey, clusteringKey); !ok {
				b.mu.Unlock()
				return false, nil
			}
		}
		b.mu.Unlock()
		return true, nil
	}
}

type buckets struct {
	// Map of [key] => [took token].
	// Took token is start from 0.
	// Then, ([took token] + 1) access was permitted.
	m sync.Map

	mu       sync.Mutex
	expunged bool
}

func newBuckets() *buckets {
	return &buckets{
		m: sync.Map{},
	}
}

// expungedBucket means that bucket was expunged.
// After swap to this value, the bucket must be deleted from buckets.
const expungedBucket = int32(math.MinInt32)

func (b *buckets) fill() bool {
	n := int32(DefaultRate / int32(time.Second/DefaultInterval))
	empty := true
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
				empty = false
				break
			}
		}

		return true
	})
	return empty
}

func (b *buckets) empty() bool {
	empty := true
	b.m.Range(func(key, value interface{}) bool {
		p := value.(*int32)

		current := atomic.LoadInt32(p)
		if current == expungedBucket {
			return true
		}

		empty = false
		return false

	})
	return empty
}

func (b *buckets) take(partitionKey, clusteringKey string) bool {
	newValue := new(int32)
LOAD_OR_NEW_LOOP:
	for {
		value, loaded := b.loadOrStore(clusteringKey, newValue)
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

			if atomic.CompareAndSwapInt32(value, current, next) {
				return true
			}
		}
	}
}

func (b *buckets) loadOrStore(clusteringKey string, newValue *int32) (*int32, bool) {
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
