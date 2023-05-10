package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/shumkovdenis/services/integration/helpers"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LaunchInput struct {
	PlayerId string `json:"playerId" validate:"required"`
	GameId   string `json:"gameId" validate:"required"`
}

type LaunchOutput struct {
	Token string `json:"token"`
}

func NewHTTPServer(cfg Config) error {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(
		helpers.NewTraceMiddleware(),
		helpers.NewLoggerMiddleware(),
	)
	app.Post("/launch", Launch)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func Launch(ctx *fiber.Ctx) error {
	var in LaunchInput

	if err := ctx.BodyParser(&in); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := helpers.Validate(&in); len(err) != 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	var out LaunchOutput
	out.Token = in.PlayerId + "_" + in.GameId

	return ctx.JSON(&out)
}
