package main

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	grpcUtils "github.com/shumkovdenis/bl/grpc"
	pb "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	"google.golang.org/grpc"
)

type grpcService struct {
	pb.ExampleServiceServer
	caller Callee
}

func NewGRPCService(cfg Config, caller Callee) error {
	service := grpcService{caller: caller}

	server := grpc.NewServer(grpcUtils.WithServerOptions())

	pb.RegisterExampleServiceServer(server, &service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}

	return server.Serve(lis)
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

		return nil, err
	}

	res := &pb.CallResponse{
		Content: msg.Content,
	}

	return res, nil
}
