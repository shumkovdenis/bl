package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/caarlos0/env/v8"
	integration "github.com/shumkovdenis/bl/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/bl/gen/integration/v1/integrationv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
	Port int `env:"PORT" envDefault:"6000"`
}

type Server struct{}

func (s *Server) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: 100500,
	})
	res.Header().Set("integration-version", "v1")

	log.Println("GetBalance")

	return res, nil
}

func main() {
	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		panic(err)
	}

	server := &Server{}

	mux := http.NewServeMux()
	mux.Handle(integrationConnect.NewIntegrationServiceHandler(server))

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		panic(err)
	}
}
