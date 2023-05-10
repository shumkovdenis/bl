package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
	"github.com/shumkovdenis/services/integration/helpers"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type GRPCServer struct {
	integrationConnect.UnimplementedIntegrationServiceHandler
	integrationService integrationConnect.IntegrationServiceClient
}

func NewGRPCServer(cfg Config) error {
	integrationService := integrationConnect.NewIntegrationServiceClient(
		helpers.NewInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
		connect.WithInterceptors(
			helpers.NewTraceInterceptor(),
			helpers.NewLoggerInterceptor(),
			helpers.NewAppInterceptor("remote"),
		),
	)

	server := GRPCServer{
		integrationService: integrationService,
	}

	mux := http.NewServeMux()
	mux.Handle(integrationConnect.NewIntegrationServiceHandler(
		&server,
		connect.WithInterceptors(
			helpers.NewTraceInterceptor(),
			helpers.NewLoggerInterceptor(),
		),
	))

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
	reqBalance := connect.NewRequest(&integration.GetBalanceRequest{})

	resBalance, err := s.integrationService.GetBalance(ctx, reqBalance)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: resBalance.Msg.Balance,
	})

	return res, nil
}