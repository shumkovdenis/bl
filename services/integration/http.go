package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
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

type HTTPServer struct {
	client *req.Client
}

func NewHTTPServer(cfg Config) error {
	client := req.C().
		SetBaseURL(fmt.Sprintf("http://localhost:%d", cfg.Dapr.HTTPPort)).
		WrapRoundTripFunc(
			helpers.NewClientLoggerMiddleware(),
			helpers.NewClientTraceMiddleware(cfg.HTTPTrace),
			helpers.NewClientAppMiddleware("remote"),
		)

	server := HTTPServer{client: client}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(
		helpers.NewServerTraceMiddleware(),
		helpers.NewServerLoggerMiddleware(),
	)
	app.Post("/http", server.Integration)
	app.Post("/launch", Launch)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (s *HTTPServer) Integration(ctx *fiber.Ctx) error {
	data := fiber.Map{}

	res, err := s.client.R().
		SetContext(ctx.UserContext()).
		SetSuccessResult(&data).
		Post("/http")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if !res.IsSuccessState() {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": res.String(),
		})
	}

	return ctx.JSON(&data)
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
