package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

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
	app.Use(logger.New(logger.Config{
		CustomTags: map[string]logger.LogFunc{
			"traceparent": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString(c.Get("traceparent"))
			},
			"tracestate": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString(c.Get("tracestate"))
			},
			"grpc-trace-bin": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString(c.Get("grpc-trace-bin"))
			},
		},
	}))
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

	if err := validate.Struct(&in); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var out LaunchOutput
	out.Token = in.PlayerId + "_" + in.GameId

	return ctx.JSON(&out)
}
