package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/dapr/go-sdk/service/common"
	"github.com/rs/zerolog/log"
)

func defaultHealthCheckHandler(ctx context.Context) error {
	return nil
}

// server is the gRPC service implementation for Dapr.
type server struct {
	pb.UnimplementedAppCallbackHealthCheckServer
	healthCheckHandler common.HealthCheckHandler
	grpcServer         *grpc.Server
	started            uint32
}

// NewServer creates new Service.
func NewServer(opts ...grpc.ServerOption) *server {
	var o []grpc.ServerOption
	o = append(o, WithServerOptions())
	o = append(o, opts...)

	s := server{
		healthCheckHandler: defaultHealthCheckHandler,
	}

	gs := grpc.NewServer(o...)
	pb.RegisterAppCallbackHealthCheckServer(gs, &s)

	s.grpcServer = gs

	return &s
}

// Start registers the server and starts it.
func (s *server) Start(address string) error {
	if address == "" {
		return errors.New("empty address")
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to TCP listen on %s: %w", address, err)
	}
	if !atomic.CompareAndSwapUint32(&s.started, 0, 1) {
		return errors.New("a gRPC server can only be started once")
	}
	return s.grpcServer.Serve(lis)
}

// Stop stops the previously-started service.
func (s *server) Stop() error {
	if atomic.LoadUint32(&s.started) == 0 {
		return nil
	}
	s.grpcServer.Stop()
	s.grpcServer = nil
	return nil
}

// GrecefulStop stops the previously-started service gracefully.
func (s *server) GracefulStop() error {
	if atomic.LoadUint32(&s.started) == 0 {
		return nil
	}
	s.grpcServer.GracefulStop()
	s.grpcServer = nil
	return nil
}

// GrpcServer returns the grpc.Server object managed by the server.
func (s *server) GrpcServer() *grpc.Server {
	return s.grpcServer
}

func (s *server) Run(address string) {
	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// run blocking call in a separate goroutine, report errors via channel
	go func() {
		if err := s.Start(address); err != nil {
			errChan <- err
		}
	}()

	// terminate your environment gracefully before leaving main function
	defer func() {
		log.Info().Msg("stopping server")
		s.GracefulStop()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Fatal().Err(err).Msg("server failed")
	case <-stopChan:
	}
}

// AddHealthCheckHandler appends provided app health check handler.
func (s *server) AddHealthCheckHandler(_ string, fn common.HealthCheckHandler) error {
	if fn == nil {
		return fmt.Errorf("health check handler required")
	}

	s.healthCheckHandler = fn

	return nil
}

// HealthCheck check app health status.
func (s *server) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*pb.HealthCheckResponse, error) {
	if s.healthCheckHandler != nil {
		if err := s.healthCheckHandler(ctx); err != nil {
			return &pb.HealthCheckResponse{}, err
		}

		return &pb.HealthCheckResponse{}, nil
	}

	return nil, fmt.Errorf("health check handler not implemented")
}
