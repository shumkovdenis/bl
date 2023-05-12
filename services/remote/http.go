package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shumkovdenis/services/remote/helpers"
)

func NewHTTPServer(cfg Config) error {
	app := fiber.New()
	app.Use(
		helpers.NewServerTraceMiddleware(),
		helpers.NewServerLoggerMiddleware(),
	)
	app.Post("/remote", remoteHandler)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func remoteHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"msg": "remote"})
}
