// Package client provides a client to access to token bucket server.
package client

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	"github.com/orgplace/trafficquota/proto"
	"google.golang.org/grpc"
)

// Client is a client to access to token bucket server.
type Client interface {
	// Take takes tokens from server.
	Take(chunkKey string, bucketKeys ...string) (bool, error)
	// TakeContext takes tokens from server.
	TakeContext(ctx context.Context, chunkKey string, bucketKeys ...string) (bool, error)
	// Ping verifies a connection to the server is still alive.
	Ping() error
	// PingContext verifies a connection to the server is still alive.
	PingContext(ctx context.Context) error
	// Close closes a connection.
	Close() error
}

var errNotServing = errors.New("server is not serving")

type client struct {
	clientConn   *grpc.ClientConn
	trafficQuota proto.TrafficQuotaClient
	health       grpc_health_v1.HealthClient
}

// NewInsecureClient constructs a new client without TLS.
func NewInsecureClient(addr string) (Client, error) {
	addr, options := parseAddr(addr)
	cc, err := grpc.Dial(addr, append(
		options,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithInsecure(),
	)...)
	if err != nil {
		return nil, err
	}

	return newClient(cc), nil
}

func newClient(cc *grpc.ClientConn) *client {
	return &client{
		clientConn:   cc,
		trafficQuota: proto.NewTrafficQuotaClient(cc),
		health:       grpc_health_v1.NewHealthClient(cc),
	}
}

var unixDomainSocketDialer = grpc.WithDialer(func(a string, t time.Duration) (net.Conn, error) {
	return net.Dial("unix", a)
})

func parseAddr(addr string) (string, []grpc.DialOption) {
	options := []grpc.DialOption{}

	const unixDomainSocketPrefix = "unix:"
	if strings.HasPrefix(addr, unixDomainSocketPrefix) {
		addr = addr[len(unixDomainSocketPrefix):]
		options = append(options, unixDomainSocketDialer)
	}

	return addr, options
}

func (c *client) Take(chunkKey string, bucketKeys ...string) (bool, error) {
	return c.TakeContext(context.Background(), chunkKey, bucketKeys...)
}

func (c *client) TakeContext(ctx context.Context, chunkKey string, bucketKeys ...string) (bool, error) {
	res, err := c.trafficQuota.Take(ctx, &proto.TakeRequest{
		ChunkKey:   chunkKey,
		BucketKeys: bucketKeys,
	})
	if err != nil {
		return false, err
	}
	return res.Allowed, err
}

func (c *client) Ping() error {
	return c.PingContext(context.Background())
}

func (c *client) PingContext(ctx context.Context) error {
	res, err := c.health.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if res.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return errNotServing
	}
	return nil
}

func (c *client) Close() error {
	return c.clientConn.Close()
}
