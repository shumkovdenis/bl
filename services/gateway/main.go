package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/caarlos0/env/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/metadata"

	"github.com/bufbuild/connect-go"
	integration "github.com/shumkovdenis/bl/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/bl/gen/integration/v1/integrationv1connect"
)

type Config struct {
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

	// client, err := dapr.NewClient()
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Close()

	// sessionService := services.NewDaprSessionStore(client, "statestore")
	// betService := services.NewDaprBetService(client)
	// historyService := services.NewDaprHistoryService(client, "history-pubsub", "history")
	// command := services.NewCommandService(sessionService, betService, historyService)

	client := integrationConnect.NewIntegrationServiceClient(
		newInsecureClient(),
		"http://localhost:50001",
		connect.WithGRPC(),
	)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		ctx := metadata.AppendToOutgoingContext(context.Background(), "dapr-app-id", "integration")

		req := connect.NewRequest(&integration.GetBalanceRequest{PlayerId: "1"})
		req.Header().Set("dapr-app-id", "integration")

		res, err := client.GetBalance(ctx, req)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println(res.Msg.Balance)

		return c.SendString(fmt.Sprintf("Balance: %d", res.Msg.Balance))
	})

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
