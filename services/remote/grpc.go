package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type GRPCServer struct {
	pb.UnimplementedGreeterServer
}

func NewGRPCServer(cfg Config) error {
	var grpcServer GRPCServer

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, &grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}

	return server.Serve(lis)
}

func (s *GRPCServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "grpc remote"}, nil
}
