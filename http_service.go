package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/shumkovdenis/bl/http/server"
)

type httpService struct {
	caller Callee
}

func NewHTTPService(cfg Config, caller Callee) error {
	s := httpService{caller: caller}

	server := server.NewServer()
	server.Post("/call", s.Handler)

	return server.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (s *httpService) Handler(c *fiber.Ctx) error {
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
