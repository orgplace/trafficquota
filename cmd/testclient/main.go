package main

import (
	"go.uber.org/zap"
	"fmt"
	"context"
	"net"

	"github.com/orgplace/trafficquota/proto"

	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewDevelopment()

	conn, err := grpc.Dial(net.JoinHostPort("localhost", "3895"), grpc.WithInsecure())
	if err != nil {
		logger.Panic("failed to dial", zap.Error(err))
	}
	defer conn.Close()

	c := proto.NewTrafficQuotaServiceClient(conn)

	res, err := c.TakeToken(context.Background(), &proto.TakeTokenRequest{
		PartitionKey: "sample",
		ClustringKey: []string{"test"},
	})
	if err != nil {
		logger.Panic("failed to take token", zap.Error(err))
	}

	fmt.Printf("%#v", res)
}
