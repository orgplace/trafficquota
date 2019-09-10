package tokenbucket

import (
	"testing"
	"time"
)

func TestTimeSliceTokenBucket_Take(t *testing.T) {
	t.Parallel()

	filledPerInterval := DefaultRate / int32(time.Second/DefaultTimeSlice)

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
			tb := NewTimeSliceTokenBucket(DefaultConfig)

			seq := 1
			for _, req := range tt.params.requests {
				tb.Fill()
				for i := int32(0); i < req; i++ {
					ok, _ := tb.Take("chunkKey", []string{"bucketKey"})
					if (seq == tt.params.notConformantAt) == ok {
						t.Errorf("Unexpected conformant: %d", seq)
					}
					seq++
				}
			}
		})
	}
}

func TestTimeSliceTokenBucket_Take_expunged(t *testing.T) {
	t.Parallel()

	tb := &timeSliceTokenBucket{}

	expungedChunk := newChunk()
	expungedValue := DefaultBucketSize
	expungedChunk.expunged = true
	expungedChunk.m.Store("bucketKey", expungedValue)
	tb.m.Store("chunkKey", expungedChunk)

	ok, _ := tb.Take("chunkKey", []string{"bucketKey"})
	if !ok {
		t.Error("could not take a token")
	}
}

func TestChunk_fill_expunged(t *testing.T) {
	t.Parallel()

	chunk := newChunk()

	expungedValue := int32(expungedBucket)
	chunk.m.Store("clustringKey", &expungedValue)

	if !chunk.fill(DefaultConfig, "chunkKey") {
		t.Error("chunk must be empty")
	}
}

func TestChunk_empty_not_empty(t *testing.T) {
	t.Parallel()

	chunk := newChunk()

	chunk.m.Store("clustringKey", new(int32))

	if chunk.empty() {
		t.Error("chunk must not be empty")
	}
}

func TestChunk_empty_expunged(t *testing.T) {
	t.Parallel()

	chunk := newChunk()

	expungedValue := int32(expungedBucket)
	chunk.m.Store("clustringKey", &expungedValue)

	if !chunk.empty() {
		t.Error("chunk must be empty")
	}
}

func TestChunk_take_expunged(t *testing.T) {
	t.Parallel()

	chunk := newChunk()

	expungedValue := int32(expungedBucket)
	chunk.m.Store("clustringKey", &expungedValue)

	if !chunk.take(DefaultConfig, "", "clustringKey") {
		t.Error("token must be taken")
	}
}
