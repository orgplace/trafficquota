package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/orgplace/trafficquota/config"
	"github.com/orgplace/trafficquota/proto"
	"github.com/orgplace/trafficquota/tokenbucket"

	"google.golang.org/grpc"
)

func newConnection(addr string) (*grpc.ClientConn, error) {
	const unixSocketPrefix = "unix:"

	if strings.HasPrefix(addr, unixSocketPrefix) {
		socketFile := addr[len(unixSocketPrefix):]

		return grpc.Dial(
			socketFile,
			grpc.WithInsecure(),
			grpc.WithDialer(func(a string, t time.Duration) (net.Conn, error) {
				return net.Dial("unix", a)
			}),
			// grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(logger)),
		)
	}

	return grpc.Dial(addr, grpc.WithInsecure())
}

func main() {
	logger, _ := zap.NewDevelopment()

	conn, err := newConnection(config.Listen)
	if err != nil {
		logger.Panic("failed to dial", zap.Error(err))
	}
	defer conn.Close()

	c := proto.NewTrafficQuotaServiceClient(conn)

	for i := 0; i < tokenbucket.DefaultBucketSize; i++ {

		res, err := c.TakeToken(context.Background(), &proto.TakeTokenRequest{
			PartitionKey:  "sample",
			ClusteringKey: []string{"test"},
		})
		if err != nil {
			logger.Panic("failed to take token", zap.Error(err))
		}

		fmt.Printf("%v\n", res.Allowed)
	}
}
