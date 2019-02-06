package server

import (
	"context"

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

func (s *trafficQuotaServer) Take(ctx context.Context, req *proto.TakeRequest) (*proto.TakeResponse, error) {
	ok, err := s.tokenBucket.Take(req.PartitionKey, req.ClusteringKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.TakeResponse{Allowed: ok}, nil
}
