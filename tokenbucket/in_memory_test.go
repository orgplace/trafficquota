package tokenbucket

import (
	"testing"
	"time"
)

func TestTake(t *testing.T) {
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
