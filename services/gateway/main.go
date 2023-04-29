package main

import (
	"fmt"

	"github.com/caarlos0/env/v8"
	"github.com/gofiber/fiber/v2"
)

type config struct {
	Port int `env:"PORT" envDefault:"6000"`
}

func main() {
	cfg := config{}
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

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
