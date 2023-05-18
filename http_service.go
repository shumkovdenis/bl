package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	httpUtils "github.com/shumkovdenis/bl/http"
)

type httpService struct {
	caller Callee
}

func NewHTTPService(cfg Config, caller Callee) error {
	service := httpService{caller: caller}

	server := httpUtils.NewServer()
	server.Post("/call", service.Handler)

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// run blocking call in a separate goroutine, report errors via channel
	go func() {
		if err := server.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			errChan <- err
		}
	}()

	// terminate your environment gracefully before leaving main function
	defer func() {
		server.Shutdown()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Fatal().Err(err).Msg("server failed")
	case <-stopChan:
	}

	return nil
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
