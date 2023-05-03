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
	integration "github.com/shumkovdenis/bl/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/bl/gen/integration/v1/integrationv1connect"
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
	// ServiceName string            `env:"SERVICE_NAME" envDefault:"gateway"`
	Port int `env:"PORT" envDefault:"6000"`
}

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}

func main() {
	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		panic(err)
	}

	client := integrationConnect.NewIntegrationServiceClient(
		newInsecureClient(),
		fmt.Sprintf("localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
	)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		req := connect.NewRequest(&integration.GetBalanceRequest{PlayerId: "1"})
		req.Header().Set("dapr-app-id", cfg.Integration.AppID)

		res, err := client.GetBalance(context.Background(), req)
		if err != nil {
			log.Println(connect.CodeOf(err))
			if connectErr := new(connect.Error); errors.As(err, &connectErr) {
				log.Println(connectErr.Message())
				log.Println(connectErr.Details())
			}
			return err
		}

		return c.SendString(fmt.Sprintf("Balance: %d", res.Msg.Balance))
	})

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
