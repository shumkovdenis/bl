package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/dapr/go-sdk/service/common"
)

func defaultHealthCheckHandler(ctx context.Context) error {
	return nil
}

// server is the gRPC service implementation for Dapr.
type server struct {
	pb.UnimplementedAppCallbackHealthCheckServer
	listener           net.Listener
	healthCheckHandler common.HealthCheckHandler
	grpcServer         *grpc.Server
	started            uint32
}

func newServer(lis net.Listener, opts ...grpc.ServerOption) *server {
	s := &server{
		listener:           lis,
		healthCheckHandler: defaultHealthCheckHandler,
	}

	gs := grpc.NewServer(opts...)
	pb.RegisterAppCallbackHealthCheckServer(gs, s)
	s.grpcServer = gs

	return s
}

// NewServer creates new Service.
func NewServer(address string, opts ...grpc.ServerOption) (*server, error) {
	if address == "" {
		return nil, errors.New("empty address")
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to TCP listen on %s: %w", address, err)
	}

	var o []grpc.ServerOption
	o = append(o, WithServerOptions())
	o = append(o, opts...)

	s := newServer(lis, o...)

	return s, nil
}

// Start registers the server and starts it.
func (s *server) Start() error {
	if !atomic.CompareAndSwapUint32(&s.started, 0, 1) {
		return errors.New("a gRPC server can only be started once")
	}
	return s.grpcServer.Serve(s.listener)
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
