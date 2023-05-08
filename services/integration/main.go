package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/caarlos0/env/v8"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

type BalanceData struct {
	Balance int64 `json:"balance"`
}

type DaprConfig struct {
	HHTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type WalletConfig struct {
	BindingName string `env:"BINDING_NAME" envDefault:"wallet"`
}

type Config struct {
	Dapr   DaprConfig   `envPrefix:"DAPR_"`
	Wallet WalletConfig `envPrefix:"WALLET_"`
	Port   int          `env:"PORT" envDefault:"6000"`
}

type Server struct {
	walletBindingName string
	client            integrationConnect.IntegrationServiceClient
}

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}

func newError(transactionID string) error {
	err := connect.NewError(
		connect.CodeInvalidArgument,
		errors.New("player id is required"),
	)

	rollbackInfo := &integration.RollbackInfo{
		TransactionId: transactionID,
	}

	if detail, detailErr := connect.NewErrorDetail(rollbackInfo); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func (s *Server) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, err // automatically coded correctly
	}

	if req.Msg.PlayerId == "" {
		return nil, newError("123")
	}

	log.Println("req:traceparent", req.Header().Get("traceparent"))
	log.Println("req:tracestate", req.Header().Get("tracestate"))
	log.Println("req:grpc-trace-bin", req.Header().Get("grpc-trace-bin"))

	// client, err := dapr.NewClient()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// in := &dapr.InvokeBindingRequest{
	// 	Name:      s.walletBindingName,
	// 	Operation: "post",
	// 	Data:      []byte(""),
	// 	Metadata:  map[string]string{"path": "/6b9663d1-41a3-47f8-8e56-8e5c8678bcde"},
	// }

	// event, err := client.InvokeBinding(ctx, in)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, connect.NewError(
	// 		connect.CodeInvalidArgument,
	// 		err,
	// 	)
	// }

	// log.Println(event.Metadata)
	// log.Println(string(event.Data))
	// log.Println(event.Metadata["statusCode"])

	// data := BalanceData{}
	// if err := json.Unmarshal(event.Data, &data); err != nil {
	// 	return nil, err
	// }

	// traceparent, err := tracing.Parse(req.Header().Get("traceparent"))
	// if err != nil {
	// 	return nil, err
	// }

	// span, err := traceparent.NewSpan()
	// if err != nil {
	// 	return nil, err
	// }

	r := connect.NewRequest(&integration.GetBalanceRequest{PlayerId: "123"})
	r.Header().Set("dapr-app-id", "remote")
	r.Header().Set("traceparent", req.Header().Get("tracestate"))
	// r.Header().Set("traceparent", span.String())
	// r.Header().Set("grpc-trace-bin", req.Header().Get("grpc-trace-bin"))

	log.Println("balance-req:traceparent", r.Header().Get("traceparent"))
	log.Println("balance-req:grpc-trace-bin", r.Header().Get("grpc-trace-bin"))

	t, err := s.client.GetBalance(ctx, r)
	if err != nil {
		log.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			log.Println(connectErr.Message())
			log.Println(connectErr.Details())
		}
		return nil, err
	}

	log.Println("balance-res:traceparent", t.Header().Get("traceparent"))
	log.Println("balance-res:grpc-trace-bin", t.Header().Get("grpc-trace-bin"))

	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: t.Msg.Balance,
	})
	// res.Header().Set("traceparent", span.String())

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

	client := integrationConnect.NewIntegrationServiceClient(
		newInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
		// connect.WithInterceptors(otelconnect.NewInterceptor()),
	)

	server := &Server{
		walletBindingName: cfg.Wallet.BindingName,
		client:            client,
	}

	mux := http.NewServeMux()
	mux.Handle(
		integrationConnect.NewIntegrationServiceHandler(
			server,
			// connect.WithInterceptors(otelconnect.NewInterceptor(otelconnect.WithTrustRemote())),
		),
	)

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}
