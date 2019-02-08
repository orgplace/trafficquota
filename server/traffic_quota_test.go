package server

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"

	"github.com/orgplace/trafficquota/tokenbucket"

	"github.com/orgplace/trafficquota/proto"
	"go.uber.org/zap/zaptest"
)

func Test_trafficQuotaServer_Take(t *testing.T) {
	type args struct {
		ctx context.Context
		req *proto.TakeRequest
	}
	tests := []struct {
		name               string
		args               args
		tokenBucketAllowed bool
		tokenBucketError   error
		want               *proto.TakeResponse
		wantErr            error
	}{
		{
			name: "allowed empty",
			args: args{
				req: &proto.TakeRequest{},
			},
			tokenBucketAllowed: true,
			want:               &proto.TakeResponse{Allowed: true},
		},
		{
			name: "error",
			args: args{
				req: &proto.TakeRequest{},
			},
			tokenBucketAllowed: true,
			tokenBucketError:   errors.New("error for test"),
			wantErr:            status.Error(codes.Internal, "error for test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenBucket := tokenbucket.NewMockTokenBucket(ctrl)
			tokenBucket.EXPECT().Take(
				tt.args.req.PartitionKey,
				tt.args.req.ClusteringKeys,
			).Return(tt.tokenBucketAllowed, tt.tokenBucketError)

			s := NewTrafficQuotaServer(logger, tokenBucket)
			got, err := s.Take(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trafficQuotaServer.Take() = %v, want %v", got, tt.want)
			}
		})
	}
}
