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

type ConnectServer struct {
	integrationConnect.UnimplementedIntegrationServiceHandler
	integrationService integrationConnect.IntegrationServiceClient
}

func NewConnectServer(cfg Config) error {
	integrationService := integrationConnect.NewIntegrationServiceClient(
		helpers.NewInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
		connect.WithInterceptors(
			helpers.NewAppInterceptor("remote"),
			helpers.NewTraceInterceptor(cfg.GRPCTrace),
			helpers.NewLoggerInterceptor(),
		),
	)

	server := ConnectServer{
		integrationService: integrationService,
	}

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
