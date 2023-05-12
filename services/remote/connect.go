package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
	"github.com/shumkovdenis/services/remote/helpers"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ConnectServer struct {
	integrationConnect.UnimplementedIntegrationServiceHandler
}

func NewConnectServer(cfg Config) error {
	var server ConnectServer

	mux := http.NewServeMux()
	mux.Handle(integrationConnect.NewIntegrationServiceHandler(
		&server,
		connect.WithInterceptors(
			helpers.NewTraceInterceptor(cfg.GRPCTrace),
			helpers.NewLoggerInterceptor(),
		),
	))

	return http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func (s *ConnectServer) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: 9999,
	})
	return res, nil
}
