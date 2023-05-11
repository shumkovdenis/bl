package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shumkovdenis/bl/services/gateway/helpers"
)

func NewHTTPServer(cfg Config) error {
	app := fiber.New()
	app.Use(
		helpers.NewTraceMiddleware(),
		helpers.NewLoggerMiddleware(),
	)
	app.Post("/remote", remoteHandler)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func remoteHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"msg": "remote"})
}
