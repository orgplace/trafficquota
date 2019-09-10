package tokenbucket

import (
	"testing"
)

func TestBucketOption_getSize(t *testing.T) {
	tests := []struct {
		name string
		o    BucketOption
		def  int32
		want int32
	}{
		{
			name: "",
			o: BucketOption{
				Size: 2,
			},
			def:  1,
			want: 2,
		},
		{
			name: "default",
			def:  1,
			want: 1,
		},
		{
			name: "banned",
			o: BucketOption{
				Banned: true,
				Size:   2,
			},
			def:  1,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.getSize(tt.def); got != tt.want {
				t.Errorf("BucketOption.getSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketOption_getRate(t *testing.T) {
	tests := []struct {
		name string
		o    BucketOption
		def  int32
		want int32
	}{
		{
			name: "",
			o: BucketOption{
				Rate: 2,
			},
			def:  1,
			want: 2,
		},
		{
			name: "default",
			def:  1,
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.getRate(tt.def); got != tt.want {
				t.Errorf("BucketOption.getRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
