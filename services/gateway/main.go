package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/caarlos0/env/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/http2"

	"github.com/bufbuild/connect-go"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
)

type DaprConfig struct {
	HHTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type IntegrationConfig struct {
	AppID string `env:"APP_ID" envDefault:"integration"`
}

type Config struct {
	Dapr        DaprConfig        `envPrefix:"DAPR_"`
	Integration IntegrationConfig `envPrefix:"INTEGRATION_"`
	Port        int               `env:"PORT" envDefault:"6000"`
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

func extractError(err error) (*integration.RollbackInfo, bool) {
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		return nil, false
	}

	for _, detail := range connectErr.Details() {
		msg, valueErr := detail.Value()
		if valueErr != nil {
			// Usually, errors here mean that we don't have the schema for this
			// Protobuf message.
			continue
		}
		if retryInfo, ok := msg.(*integration.RollbackInfo); ok {
			return retryInfo, true
		}
	}

	return nil, false
}

func main() {
	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		panic(err)
	}

	client := integrationConnect.NewIntegrationServiceClient(
		newInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
	)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		playerId := c.Query("player_id")

		log.Println("traceparent", c.Get("traceparent"))
		log.Println("grpc-trace-bin", c.Get("grpc-trace-bin"))

		req := connect.NewRequest(&integration.GetBalanceRequest{PlayerId: playerId})
		req.Header().Set("dapr-app-id", cfg.Integration.AppID)

		res, err := client.GetBalance(context.Background(), req)
		if err != nil {
			log.Println(connect.CodeOf(err))
			if connectErr := new(connect.Error); errors.As(err, &connectErr) {
				log.Println(connectErr.Message())
				log.Println(connectErr.Details())

				if rollbackInfo, ok := extractError(err); ok {
					log.Println(rollbackInfo.TransactionId)
				}
			}
			return err
		}

		return c.SendString(fmt.Sprintf("Balance: %d", res.Msg.Balance))
	})

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
