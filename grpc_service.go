package main

import (
	"context"
	"fmt"

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

func RunGRPCService(cfg Config, caller Callee) {
	service := grpcService{caller: caller}

	server := grpcUtils.NewServer()

	pb.RegisterExampleServiceServer(server.GrpcServer(), &service)

	server.Run(fmt.Sprintf(":%d", cfg.Port))
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
