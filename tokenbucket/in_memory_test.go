package tokenbucket

import (
	"testing"
	"time"
)

func TestInMemoryTokenBucket_Take(t *testing.T) {
	t.Parallel()

	filledPerInterval := DefaultRate / int32(time.Second/DefaultInterval)

	type params struct {
		requests        []int32
		notConformantAt int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			"burst",
			params{
				requests:        []int32{DefaultBucketSize + 1},
				notConformantAt: int(DefaultBucketSize) + 1,
			},
		},
		{
			"rate",
			params{
				requests:        []int32{DefaultBucketSize, filledPerInterval + 1},
				notConformantAt: int(DefaultBucketSize+filledPerInterval) + 1,
			},
		},
		{
			"fully filled",
			params{
				requests:        []int32{1, DefaultBucketSize + 1},
				notConformantAt: int(DefaultBucketSize) + 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tb := NewInMemoryTokenBucket(DefaultConfig)

			seq := 1
			for _, req := range tt.params.requests {
				tb.Fill()
				for i := int32(0); i < req; i++ {
					ok, _ := tb.Take("partitionKey", []string{"clusteringKey"})
					if (seq == tt.params.notConformantAt) == ok {
						t.Errorf("Unexpected conformant: %d", seq)
					}
					seq++
				}
			}
		})
	}
}

func TestInMemoryTokenBucket_Take_expunged(t *testing.T) {
	t.Parallel()

	tb := &inMemoryTokenBucket{}

	expungedBuckets := newBuckets()
	expungedValue := DefaultBucketSize
	expungedBuckets.expunged = true
	expungedBuckets.m.Store("clusteringKey", expungedValue)
	tb.m.Store("partitionKey", expungedBuckets)

	ok, _ := tb.Take("partitionKey", []string{"clusteringKey"})
	if !ok {
		t.Error("could not take a token")
	}
}

func TestBuckets_fill_expunged(t *testing.T) {
	t.Parallel()

	buckets := newBuckets()

	expungedValue := int32(expungedBucket)
	buckets.m.Store("clustringKey", &expungedValue)

	if !buckets.fill(DefaultConfig, "partitionKey") {
		t.Error("buckets must be empty")
	}
}

func TestBuckets_empty_not_empty(t *testing.T) {
	t.Parallel()

	buckets := newBuckets()

	buckets.m.Store("clustringKey", new(int32))

	if buckets.empty() {
		t.Error("buckets must not be empty")
	}
}

func TestBuckets_empty_expunged(t *testing.T) {
	t.Parallel()

	buckets := newBuckets()

	expungedValue := int32(expungedBucket)
	buckets.m.Store("clustringKey", &expungedValue)

	if !buckets.empty() {
		t.Error("buckets must be empty")
	}
}

func TestBuckets_take_expunged(t *testing.T) {
	t.Parallel()

	buckets := newBuckets()

	expungedValue := int32(expungedBucket)
	buckets.m.Store("clustringKey", &expungedValue)

	if !buckets.take(DefaultConfig, "", "clustringKey") {
		t.Error("token must be taken")
	}
}
