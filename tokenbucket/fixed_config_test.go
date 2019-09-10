package tokenbucket

import (
	"testing"
	"time"
)

func Test_fixedConfig_Overflow(t *testing.T) {
	type args struct {
		chunkKey  string
		bucketKey string
		tokens    int32
	}
	tests := []struct {
		name   string
		option Option
		args   args
		want   bool
	}{
		{
			name: "default",
			option: Option{
				Default: BucketOption{Size: 1},
			},
			args: args{
				tokens: 1,
			},
		},
		{
			name: "over default",
			option: Option{
				Default: BucketOption{Size: 1},
			},
			args: args{
				tokens: int32(2),
			},
			want: true,
		},
		{
			name: "chunk default",
			option: Option{
				Default: BucketOption{Size: 1},
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Default: BucketOption{Size: 1},
					},
				},
			},
			args: args{
				chunkKey: "chunkKey",
				tokens:   1,
			},
		},
		{
			name: "over chunk default",
			option: Option{
				Default: BucketOption{Size: 2},
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Default: BucketOption{Size: 1},
					},
				},
			},
			args: args{
				chunkKey: "chunkKey",
				tokens:   2,
			},
			want: true,
		},
		{
			name: "in bucket",
			option: Option{
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Buckets: map[string]*BucketOption{
							"bucketKey":  {Size: 1},
							"bucketKey2": {Size: 2},
						},
					},
				},
			},
			args: args{
				chunkKey:  "chunkKey",
				bucketKey: "bucketKey2",
				tokens:    2,
			},
		},
		{
			name: "over bucket",
			option: Option{
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Buckets: map[string]*BucketOption{
							"bucketKey": {Size: 1},
						},
					},
				},
			},
			args: args{
				chunkKey:  "chunkKey",
				bucketKey: "bucketKey",
				tokens:    1,
			},
			want: false,
		},
		{
			name: "min bucket optimization",
			option: Option{
				Default: BucketOption{Size: 1},
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Default: BucketOption{Size: 2},
					},
				},
			},
			args: args{
				chunkKey: "chunkKey",
				tokens:   2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFixedConfig(&tt.option)
			if got := c.Overflow(tt.args.chunkKey, tt.args.bucketKey, tt.args.tokens); got != tt.want {
				t.Errorf("fixedConfig.Overflow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixedConfig_Rate(t *testing.T) {
	const onePerInterval = int32(time.Second / DefaultTimeSlice)

	type args struct {
		chunkKey  string
		bucketKey string
	}
	tests := []struct {
		name   string
		option Option
		args   args
		want   int32
	}{
		{
			name: "default",
			option: Option{
				Default: BucketOption{Rate: onePerInterval * 2},
			},
			want: 2,
		},
		{
			name: "chunk default",
			option: Option{
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Default: BucketOption{Rate: onePerInterval * 2},
					},
				},
			},
			args: args{
				chunkKey: "chunkKey",
			},
			want: 2,
		},
		{
			name: "bucket rate",
			option: Option{
				Chunks: map[string]*ChunkOption{
					"chunkKey": {
						Buckets: map[string]*BucketOption{
							"bucketKey": {Rate: onePerInterval},
						},
					},
				},
			},
			args: args{
				chunkKey:  "chunkKey",
				bucketKey: "bucketKey",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFixedConfig(&tt.option)
			if got := c.Rate(tt.args.chunkKey, tt.args.bucketKey); got != tt.want {
				t.Errorf("fixedConfig.Rate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFixedConfig_default(t *testing.T) {
	if got := NewFixedConfig(&Option{}); got != DefaultConfig {
		t.Errorf("NewFixedConfig() = %v", got)
	}
}
