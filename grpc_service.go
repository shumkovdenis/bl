package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	grpcUtils "github.com/shumkovdenis/bl/grpc"
	pb "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	pb.ExampleServiceServer
	caller Callee
}

func NewGRPCService(cfg Config, caller Callee) error {
	service := grpcService{caller: caller}

	server, err := grpcUtils.NewServer(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}

	pb.RegisterExampleServiceServer(server.GrpcServer(), &service)

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// run blocking call in a separate goroutine, report errors via channel
	go func() {
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// terminate your environment gracefully before leaving main function
	defer func() {
		server.GracefulStop()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Fatal().Err(err).Msg("server failed")
	case <-stopChan:
	}

	return nil
}

func (s *grpcService) Call(ctx context.Context, req *pb.CallRequest) (*pb.CallResponse, error) {
	logger := log.Ctx(ctx)

	logger.Info().
		Str("message", req.Message).
		Msg("service received message")

	msg := Message{Content: req.Message}

	msg, err := s.caller.Call(ctx, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("service failed to call")

		st := status.New(codes.Internal, err.Error())

		ds, err := st.WithDetails(&pb.RetryInfo{
			Count:   1,
			Timeout: 100,
		})
		if err != nil {
			return nil, st.Err()
		}

		return nil, ds.Err()
	}

	res := &pb.CallResponse{
		Content: msg.Content,
	}

	return res, nil
}
