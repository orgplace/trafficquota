package tokenbucket

import (
	"math"
	"sync"
	"sync/atomic"
)

type inMemoryTokenBucket struct {
	config Config

	m sync.Map
}

// NewInMemoryTokenBucket constructs a in-memory TokenBucket
func NewInMemoryTokenBucket(config Config) TokenBucket {
	return &inMemoryTokenBucket{config: config}
}

func (tb *inMemoryTokenBucket) Fill() {
	tb.m.Range(func(key, value interface{}) bool {
		b := value.(*chunk)

		if b.fill(tb.config, key.(string)) {
			b.mu.Lock()
			if b.empty() {
				b.expunged = true
			}
			b.mu.Unlock()
			tb.m.Delete(key)
		}

		return true
	})
}

func (tb *inMemoryTokenBucket) Take(chunkKey string, bucketKeys []string) (bool, error) {
	newValue := newChunk()
	for {
		value, _ := tb.m.LoadOrStore(chunkKey, newValue)
		b := value.(*chunk)

		b.mu.RLock()
		if b.expunged {
			b.mu.RUnlock()
			tb.m.Delete(chunkKey)
			continue
		}

		for _, bucketKey := range bucketKeys {
			if ok := b.take(tb.config, chunkKey, bucketKey); !ok {
				b.mu.RUnlock()
				return false, nil
			}
		}
		b.mu.RUnlock()
		return true, nil
	}
}

type chunk struct {
	// Map of [key] => [took token].
	// Took token is start from 0.
	// Then, ([took token] + 1) access was permitted.
	m sync.Map

	mu       sync.RWMutex
	expunged bool
}

func newChunk() *chunk {
	return &chunk{}
}

// expungedBucket means that bucket was expunged.
// After swap to this value, the bucket must be deleted from chunk.
const expungedBucket = int32(math.MinInt32)

func (b *chunk) fill(config Config, chunkKey string) bool {
	empty := true
	b.m.Range(func(key, value interface{}) bool {
		p := value.(*int32)

		for {
			current := atomic.LoadInt32(p)
			if current == expungedBucket {
				b.m.Delete(key)
				break
			}

			next := current - config.Rate(chunkKey, key.(string))
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

func (b *chunk) empty() bool {
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

func (b *chunk) take(config Config, chunkKey, bucketKey string) bool {
	newValue := new(int32) // Starts from 0
LOAD_OR_NEW_LOOP:
	for {
		value, loaded := b.loadOrStore(bucketKey, newValue)
		if !loaded {
			return true
		}

		for {
			current := atomic.LoadInt32(value)
			if current == expungedBucket {
				b.m.Delete(bucketKey)
				continue LOAD_OR_NEW_LOOP
			}

			next := current + 1
			if config.Overflow(chunkKey, bucketKey, next+1) {
				return false
			}

			if atomic.CompareAndSwapInt32(value, current, next) {
				return true
			}
		}
	}
}

func (b *chunk) loadOrStore(bucketKey string, newValue *int32) (*int32, bool) {
	for {
		p, loaded := b.m.LoadOrStore(bucketKey, newValue)
		if !loaded {
			return newValue, false
		}

		value := p.(*int32)
		if atomic.LoadInt32(value) == expungedBucket {
			b.m.Delete(bucketKey)
		} else {
			return value, true
		}
	}
}
