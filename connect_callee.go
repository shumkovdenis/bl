package main

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	connectUtils "github.com/shumkovdenis/bl/connect"
	example "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	exampleConnect "github.com/shumkovdenis/protobuf-schema/gen/example/v1/examplev1connect"
)

type connectCallee struct {
	client exampleConnect.IntegrationServiceClient
}

func NewConnectCallee(cfg Config) *connectCallee {
	client := exampleConnect.NewIntegrationServiceClient(
		connectUtils.NewInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connectUtils.WithClientOptions(
			connectUtils.AddDaprAppIDHeader(cfg.Callee.ServiceName),
		),
	)
	return &connectCallee{client: client}
}

func (c *connectCallee) Call(ctx context.Context, msg Message) (Message, error) {
	req := connect.NewRequest(&example.CallRequest{
		Message: msg.Content,
	})

	res, err := c.client.Call(ctx, req)
	if err != nil {
		return Message{}, err
	}

	m := Message{Content: res.Msg.Content}

	return m, nil
}
