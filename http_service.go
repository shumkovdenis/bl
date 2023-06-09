package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	httpUtils "github.com/shumkovdenis/bl/http"
)

type httpService struct {
	caller Callee
}

func RunHTTPService(cfg Config, caller Callee) {
	service := httpService{caller: caller}

	server := httpUtils.NewServer()
	server.Post("/call", service.Handler)
	server.Run(fmt.Sprintf(":%d", cfg.Port))
}

func (s httpService) Handler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	logger := log.Ctx(ctx)

	var msg Message

	if err := c.BodyParser(&msg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	logger.Info().
		Str("content", msg.Content).
		Msg("service received message")

	msg, err := s.caller.Call(ctx, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to call callee service")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(&msg)
}
