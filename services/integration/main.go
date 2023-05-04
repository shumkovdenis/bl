package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/caarlos0/env/v8"
	integration "github.com/shumkovdenis/bl/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/bl/gen/integration/v1/integrationv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	dapr "github.com/dapr/go-sdk/client"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type BalanceData struct {
	Balance int64 `json:"balance"`
}

type WalletConfig struct {
	BindingName string `env:"BINDING_NAME" envDefault:"wallet"`
}

type Config struct {
	Wallet WalletConfig `envPrefix:"WALLET_"`
	Port   int          `env:"PORT" envDefault:"6000"`
}

type Server struct {
	client            dapr.Client
	walletBindingName string
}

func (s *Server) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, err // automatically coded correctly
	}

	if req.Msg.PlayerId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("player_id is required"))
	}

	in := &dapr.InvokeBindingRequest{
		Name:      s.walletBindingName,
		Operation: "post",
		Data:      []byte(""),
		Metadata:  map[string]string{"path": "/6b9663d1-41a3-47f8-8e56-8e5c8678bcde"},
	}

	event, err := s.client.InvokeBinding(ctx, in)
	if err != nil {
		return nil, err
	}

	log.Println(string(event.Data))

	data := BalanceData{}
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return nil, err
	}

	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: data.Balance,
	})

	return res, nil
}

func main() {
	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatal(err)
	}

	client, err := dapr.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	server := &Server{
		client:            client,
		walletBindingName: cfg.Wallet.BindingName,
	}

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
