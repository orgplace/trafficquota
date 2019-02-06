package tokenbucket

import (
	"sync"
	"sync/atomic"
	"time"
)

type inMemoryTokenBucket struct {
	shards bucketShards
}

func NewInMemoryTokenBucket() TokenBucket {
	b := &inMemoryTokenBucket{
		shards: *newBucketShard(),
	}

	go func() {
		c := time.Tick(DefaultInterval)
		for range c {
			b.shards.each(func(partitionKey string, b *buckets) {
				b.fill()
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

func (b *buckets) fill() {
	n := int32(DefaultRate / int32(time.Second/DefaultInterval))
	b.m.Range(func(key, value interface{}) bool {
		p := value.(*int32)

		// LoadInt32 から Delete の間で rate 以上のアクセスが来ないことを仮定している。
		// もし LoadInt32 の後に take によって rate 以上増加すると、
		// Delete によって増分を失って過剰に許可してしまう。
		// エントリー単位でロックすれば回避できるが行っていない。
		if atomic.LoadInt32(p) < 0 {
			b.m.Delete(key)
		} else {
			atomic.AddInt32(p, -n)
		}

		return false
	})
}

func (b *buckets) take(partitionKey, clusteringKey string) bool {
	value, loaded := b.m.LoadOrStore(clusteringKey, new(int32))
	if loaded {
		p := value.(*int32)
		for {
			current := atomic.LoadInt32(p)
			if DefaultBucketSize <= current+1 {
				// TODO: load configured size
				//if configuredBucketSize <= current {
				return false
				//}
			}

			if atomic.CompareAndSwapInt32(p, current, maxInt32(current+1, 0)) {
				return true
			}
		}
	}
	return true
}

func maxInt32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}
