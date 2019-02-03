package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/orgplace/trafficquota/proto"
	"github.com/orgplace/trafficquota/tokenbucket"
	"go.uber.org/zap"
)

type trafficQuotaServer struct {
	logger      *zap.Logger
	tokenBucket tokenbucket.TokenBucket
}

func NewTrafficQuotaServer(logger *zap.Logger) proto.TrafficQuotaServiceServer {
	return &trafficQuotaServer{
		logger:      logger,
		tokenBucket: tokenbucket.NewInMemoryTokenBucket(),
	}
}

func (s *trafficQuotaServer) TakeToken(ctx context.Context, req *proto.TakeTokenRequest) (*proto.TakeTokenResponse, error) {
	s.logger.Debug("take token")
	now := time.Now()
	ok, err := s.tokenBucket.Take(req.PartitionKey, req.ClusteringKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.logger.Debug("duration", zap.Duration("time", time.Since(now)))

	return &proto.TakeTokenResponse{Allowed: ok}, nil
}
