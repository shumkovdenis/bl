package main

import (
	"context"
	"fmt"

	grpcUtils "github.com/shumkovdenis/bl/grpc"
	pb "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type grpcCallee struct {
	cfg Config
}

func NewGRPCCallee(cfg Config) *grpcCallee {
	return &grpcCallee{cfg: cfg}
}

func (c grpcCallee) Call(ctx context.Context, msg Message) (Message, error) {
	var interceptor grpc.UnaryClientInterceptor
	if c.cfg.IsBinary {
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
		return Message{}, fmt.Errorf("did not connect: %w", err)
	}
	defer conn.Close()

	req := &pb.CallRequest{Message: msg.Content}

	client := pb.NewExampleServiceClient(conn)

	res, err := client.Call(ctx, req)
	if err != nil {
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *pb.RetryInfo:
				return Message{}, fmt.Errorf("please retry with count=%d", info.GetCount())
			}
		}

		return Message{}, fmt.Errorf("failed to call: %w", err)
	}

	m := Message{Content: res.Content}

	return m, nil
}
