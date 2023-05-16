package main

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	connectUtils "github.com/shumkovdenis/bl/connect"
	examplePb "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	exampleConnect "github.com/shumkovdenis/protobuf-schema/gen/example/v1/examplev1connect"
)

type connectCallee struct {
	client exampleConnect.IntegrationServiceClient
}

func NewConnectCallee(cfg Config) *connectCallee {
	var interceptor connect.UnaryInterceptorFunc
	if cfg.IsBinary() {
		interceptor = connectUtils.AddBinaryTraceContextHeader()
	} else {
		interceptor = connectUtils.AddTraceContextHeader()
	}

	client := exampleConnect.NewIntegrationServiceClient(
		connectUtils.NewInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connectUtils.WithClientOptions(
			connectUtils.AddDaprAppIDHeader(cfg.Callee.ServiceName),
			interceptor,
		),
	)
	return &connectCallee{client: client}
}

func (c connectCallee) Call(ctx context.Context, msg Message) (Message, error) {
	req := connect.NewRequest(&examplePb.CallRequest{
		Message: msg.Content,
	})

	res, err := c.client.Call(ctx, req)
	if err != nil {
		return Message{}, err
	}

	m := Message{Content: res.Msg.Content}

	return m, nil
}
