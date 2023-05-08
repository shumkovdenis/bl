package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/caarlos0/env/v8"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type DaprConfig struct {
	HHTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type Config struct {
	Dapr DaprConfig `envPrefix:"DAPR_"`
	Port int        `env:"PORT" envDefault:"6000"`
}

type Server struct {
}

func (s *Server) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, err // automatically coded correctly
	}

	log.Println("req:traceparent", req.Header().Get("traceparent"))
	log.Println("req:tracestate", req.Header().Get("tracestate"))
	log.Println("req:grpc-trace-bin", req.Header().Get("grpc-trace-bin"))

	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: 2000,
	})

	log.Println("res:traceparent", res.Header().Get("traceparent"))
	log.Println("res:grpc-trace-bin", res.Header().Get("grpc-trace-bin"))

	return res, nil
}

func main() {
	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatal(err)
	}

	server := &Server{}

	mux := http.NewServeMux()
	mux.Handle(integrationConnect.NewIntegrationServiceHandler(server))

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}
