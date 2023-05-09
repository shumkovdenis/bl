package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type InitInput struct {
	Token string `json:"token"`
}

type InitOutput struct {
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
	app.Post("/init", Init)
	app.Post("/bet", Init)
	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func Init(ctx *fiber.Ctx) error {
	var in InitInput

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

	var out InitOutput
	out.Token = in.Token

	return ctx.JSON(&out)
}
