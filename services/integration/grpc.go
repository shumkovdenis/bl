package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

type GRPCServer struct {
	pb.UnimplementedGreeterServer
	client pb.GreeterClient
}

func NewGRPCServer(cfg Config) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	client := pb.NewGreeterClient(conn)

	grpcServer := GRPCServer{client: client}

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

	log.Println("metadata from incoming context:", md)

	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", "remote")

	out, err := s.client.SayHello(ctx, &pb.HelloRequest{Name: in.GetName()})
	if err != nil {
		return nil, err
	}

	return &pb.HelloReply{Message: out.GetMessage()}, nil
}
