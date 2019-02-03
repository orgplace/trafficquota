package server

import (
	"context"

	"github.com/orgplace/trafficquota/proto"
)

type trafficQuotaServer struct{}

func NewTrafficQuotaServer() proto.TrafficQuotaServiceServer {
	return &trafficQuotaServer{}
}

func (s *trafficQuotaServer) TakeToken(context.Context, *proto.TakeTokenRequest) (*proto.TakeTokenResponse, error) {
	panic("Not implemented")
}
