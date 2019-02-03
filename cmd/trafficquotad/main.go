package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"

	"github.com/orgplace/trafficquota/proto"
	"github.com/orgplace/trafficquota/server"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewDevelopment()
	s := buildGRPCServer(logger)

	if err := listenAndServe(logger, s, net.JoinHostPort("localhost", "3895")); err != nil {
		logger.Panic("failed to start the server", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down the server")

	s.GracefulStop()
}

func buildGRPCServer(logger *zap.Logger) *grpc.Server {
	s := grpc.NewServer(buildGRPCServerOptions(logger)...)

	proto.RegisterTrafficQuotaServiceServer(s, server.NewTrafficQuotaServer())

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

func listenAndServe(logger *zap.Logger, s *grpc.Server, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	logger.Info("Starting the server", zap.String("address", address))

	go func() {
		if err := s.Serve(listener); err != nil {
			logger.Panic("failed to serve", zap.Error(err))
		}
	}()

	return nil
}
