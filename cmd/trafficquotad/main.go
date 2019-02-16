package main

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"

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

	server, err := buildGRPCServer(logger)
	if err != nil {
		logger.Panic("failed to build server", zap.Error(err))
	}

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

func buildGRPCServer(logger *zap.Logger) (*grpc.Server, error) {
	s := grpc.NewServer(buildGRPCServerOptions(logger)...)

	configFile, err := loadConfigFile(logger)
	if err != nil {
		return nil, err
	}

	tb := tokenbucket.NewInMemoryTokenBucket(
		tokenbucket.NewFixedConfig(configFile.TokenBucket.AsOption()),
	)
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

	return s, nil
}

func loadConfigFile(logger *zap.Logger) (*config.FileContent, error) {
	var result config.FileContent
	switch _, err := toml.DecodeFile(config.ConfigFilePath, &result); err.(type) {
	case *os.PathError:
		logger.Debug("default configuration is used", zap.Error(err))
		return &result, nil
	case nil:
		return &result, nil
	default:
		return nil, err
	}
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
