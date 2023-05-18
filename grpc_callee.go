package main

import (
	"context"
	"fmt"
	"log"

	grpcUtils "github.com/shumkovdenis/bl/grpc"
	pb "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcCallee struct {
	cfg Config
}

func NewGRPCCallee(cfg Config) *grpcCallee {
	return &grpcCallee{cfg: cfg}
}

func (c grpcCallee) Call(ctx context.Context, msg Message) (Message, error) {
	var interceptor grpc.UnaryClientInterceptor
	if c.cfg.IsBinary() {
		interceptor = grpcUtils.AddBinaryTraceContext()
	} else {
		interceptor = grpcUtils.AddTraceContext()
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", c.cfg.Dapr.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcUtils.WithClientOptions(
			grpcUtils.AddDaprAppID(c.cfg.Callee.ServiceName),
			interceptor,
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	req := &pb.CallRequest{Message: msg.Content}

	client := pb.NewExampleServiceClient(conn)

	res, err := client.Call(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	m := Message{Content: res.Content}

	return m, nil
}
