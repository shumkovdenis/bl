package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/shumkovdenis/services/integration/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

type GRPCServer struct {
	pb.UnimplementedGreeterServer
	config Config
}

func NewGRPCServer(cfg Config) error {
	grpcServer := GRPCServer{config: cfg}

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, &grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}

	return server.Serve(lis)
}

func (s *GRPCServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", s.config.Dapr.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	md, _ := metadata.FromIncomingContext(ctx)

	grpcTraceBin := md["grpc-trace-bin"][0]
	log.Println("metadata grpc-trace-bin:", grpcTraceBin)

	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", "remote")

	sc, ok := helpers.SpanContextFromBinary([]byte(grpcTraceBin))
	log.Println("sc:", sc.TraceID(), sc.SpanID(), "ok:", ok)
	// ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(md["grpc-trace-bin"][0]))

	grpc.SetHeader(ctx, metadata.Pairs("grpc-trace-bin", grpcTraceBin))

	req := &pb.HelloRequest{Name: in.GetName()}
	out, err := client.SayHello(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.HelloReply{Message: out.GetMessage()}, nil
}
