package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type GRPCServer struct {
	integrationConnect.UnimplementedIntegrationServiceHandler
}

func NewGRPCServer(cfg Config) error {
	var server GRPCServer

	mux := http.NewServeMux()
	mux.Handle(integrationConnect.NewIntegrationServiceHandler(&server))

	return http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func (s *GRPCServer) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: 999,
	})
	return res, nil
}
