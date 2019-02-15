// Package server implements a TrafficQuotaServer.
package server

import (
	"context"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/orgplace/trafficquota/proto"
	"github.com/orgplace/trafficquota/tokenbucket"
	"go.uber.org/zap"
)

var defaultChunkKeys = []string{""}

type trafficQuotaServer struct {
	logger      *zap.Logger
	tokenBucket tokenbucket.TokenBucket
}

// NewTrafficQuotaServer is a constructor of TrafficQuotaServer
func NewTrafficQuotaServer(logger *zap.Logger, tokenBucket tokenbucket.TokenBucket) proto.TrafficQuotaServer {
	return &trafficQuotaServer{
		logger:      logger,
		tokenBucket: tokenBucket,
	}
}

func (s *trafficQuotaServer) Take(ctx context.Context, req *proto.TakeRequest) (*proto.TakeResponse, error) {
	chunkKeys := req.ChunkKeys
	if len(chunkKeys) == 0 {
		chunkKeys = defaultChunkKeys
	}

	ok, err := s.tokenBucket.Take(req.PartitionKey, chunkKeys)
	if err != nil {
		s.logger.Error("failed to take token", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.TakeResponse{Allowed: ok}, nil
}
