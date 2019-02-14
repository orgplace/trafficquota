package tokenbucket

import (
	"testing"
	"time"
)

func TestInMemoryTokenBucket_Take(t *testing.T) {
	filledPerInterval := DefaultRate / int(time.Second/DefaultInterval)

	type params struct {
		requests        []int
		notConformantAt int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			"burst",
			params{
				requests:        []int{DefaultBucketSize + 1},
				notConformantAt: DefaultBucketSize + 1,
			},
		},
		{
			"rate",
			params{
				requests:        []int{DefaultBucketSize, filledPerInterval + 1},
				notConformantAt: DefaultBucketSize + filledPerInterval + 1,
			},
		},
		{
			"fully filled",
			params{
				requests:        []int{1, DefaultBucketSize + 1},
				notConformantAt: DefaultBucketSize + 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewInMemoryTokenBucket()

			seq := 1
			for _, req := range tt.params.requests {
				tb.Fill()
				for i := 0; i < req; i++ {
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
	buckets := newBuckets()

	expungedValue := int32(expungedBucket)
	buckets.m.Store("clustringKey", &expungedValue)

	if !buckets.fill() {
		t.Error("buckets must be empty")
	}
}

func TestBuckets_empty_not_empty(t *testing.T) {
	buckets := newBuckets()

	buckets.m.Store("clustringKey", new(int32))

	if buckets.empty() {
		t.Error("buckets must not be empty")
	}
}

func TestBuckets_empty_expunged(t *testing.T) {
	buckets := newBuckets()

	expungedValue := int32(expungedBucket)
	buckets.m.Store("clustringKey", &expungedValue)

	if !buckets.empty() {
		t.Error("buckets must be empty")
	}
}

func TestBuckets_take_expunged(t *testing.T) {
	buckets := newBuckets()

	expungedValue := int32(expungedBucket)
	buckets.m.Store("clustringKey", &expungedValue)

	if !buckets.take("", "clustringKey") {
		t.Error("token must be taken")
	}
}
