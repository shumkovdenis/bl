package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
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
	md, _ := metadata.FromIncomingContext(ctx)

	traceParent := md["traceparent"]
	grpcTraceBin := md["grpc-trace-bin"][0]
	log.Println("metadata traceParent:", traceParent)
	log.Println("metadata grpc-trace-bin:", grpcTraceBin)

	return &pb.HelloReply{Message: "grpc remote"}, nil
}
