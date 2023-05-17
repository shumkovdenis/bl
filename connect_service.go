package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	connectUtils "github.com/shumkovdenis/bl/connect"
	example "github.com/shumkovdenis/protobuf-schema/gen/example/v1"
	exampleConnect "github.com/shumkovdenis/protobuf-schema/gen/example/v1/examplev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type connectService struct {
	exampleConnect.UnimplementedIntegrationServiceHandler
	caller Callee
}

func NewConnectService(cfg Config, caller Callee) error {
	s := connectService{caller: caller}

	mux := http.NewServeMux()
	mux.Handle(exampleConnect.NewIntegrationServiceHandler(&s,
		connectUtils.WithHandlerOptions(),
	))

	return http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func (s connectService) Call(
	ctx context.Context,
	req *connect.Request[example.CallRequest],
) (*connect.Response[example.CallResponse], error) {
	logger := log.Ctx(ctx)

	logger.Info().
		Str("message", req.Msg.Message).
		Msg("service received message")

	msg := Message{Content: req.Msg.Message}

	msg, err := s.caller.Call(ctx, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("service failed to call")

		return nil, newCallError(err)
	}

	res := connect.NewResponse(&example.CallResponse{
		Content: msg.Content,
	})

	return res, nil
}
