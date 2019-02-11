package main

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/orgplace/trafficquota/config"
	"github.com/orgplace/trafficquota/server"
	"github.com/orgplace/trafficquota/tokenbucket"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"

	"github.com/orgplace/trafficquota/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

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

	server := buildGRPCServer(logger)

	if err := listenAndServe(logger, server, config.Listen); err != nil {
		logger.Panic("failed to start the server", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down the server")

	server.GracefulStop()
}

func buildGRPCServer(logger *zap.Logger) *grpc.Server {
	s := grpc.NewServer(buildGRPCServerOptions(logger)...)

	tb := tokenbucket.NewInMemoryTokenBucket()
	go func() {
		c := time.Tick(tokenbucket.DefaultInterval)
		for range c {
			tb.Fill()
		}
	}()

	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	proto.RegisterTrafficQuotaServer(s, server.NewTrafficQuotaServer(
		logger, tb,
	))

	return s
}

func buildGRPCServerOptions(logger *zap.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}
}

func listenAndServe(logger *zap.Logger, s *grpc.Server, listen string) error {
	const unixSocketPrefix = "unix:"

	var listener net.Listener

	if strings.HasPrefix(listen, unixSocketPrefix) {
		socketFile := listen[len(unixSocketPrefix):]
		os.Remove(socketFile)
		l, err := net.Listen("unix", socketFile)
		if err != nil {
			return err
		}
		listener = l

		if err := os.Chmod(socketFile, 0660); err != nil {
			return err
		}
	} else {
		l, err := net.Listen("tcp", listen)
		if err != nil {
			return err
		}
		listener = l
	}

	logger.Info("Starting the server", zap.String("listen", listen))

	go func() {
		if err := s.Serve(listener); err != nil {
			logger.Panic("failed to serve", zap.Error(err))
		}
	}()

	return nil
}
