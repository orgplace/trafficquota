package main

import (
	"context"
	"fmt"
	"net"
	"sort"
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

func newLogger() (*zap.Logger, error) {
	c := zap.NewProductionConfig()
	c.Level = zap.NewAtomicLevelAt(config.LogLevel)
	return c.Build()
}

func main() {
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Debug(config.Listen)
	conn, err := newConnection(config.Listen)
	if err != nil {
		logger.Panic("failed to dial", zap.Error(err))
	}
	defer conn.Close()

	c := proto.NewTrafficQuotaClient(conn)

	type result struct {
		allowed  bool
		duration time.Duration
	}
	n := tokenbucket.DefaultBucketSize * 5
	results := make(chan *result, n)
	for i := 0; i < n; i++ {
		go func() {
			t := time.Now()
			res, err := c.Take(context.Background(), &proto.TakeRequest{
				PartitionKey:   "sample",
				ClusteringKeys: []string{"test"},
			})
			d := time.Since(t)
			if err != nil {
				logger.Panic("failed to take token", zap.Error(err))
			}

			results <- &result{
				allowed:  res.Allowed,
				duration: d,
			}
		}()
	}

	allow := 0
	durations := make([]time.Duration, n)

	for i := 0; i < n; i++ {
		r := <-results
		if r.allowed {
			allow++
		}
		durations[i] = r.duration
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	fmt.Printf("allow: %d, deny: %d\n", allow, n-allow)

	percentile := 1
	for i, d := range durations {
		p := float64((i+1)*10) / float64(n)
		if p >= float64(percentile) {
			fmt.Printf("%6.2f%%: %s\n", p*10., d.String())
			percentile++
		}
	}
}
