package http

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type server struct {
	*fiber.App
}

func NewServer() *server {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(
		InjectTraceContext(),
		InjectTraceContextLogger(),
	)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	s := server{
		App: app,
	}

	return &s
}

func (s *server) Run(address string) {
	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// run blocking call in a separate goroutine, report errors via channel
	go func() {
		if err := s.App.Listen(address); err != nil {
			errChan <- err
		}
	}()

	// terminate your environment gracefully before leaving main function
	defer func() {
		log.Info().Msg("stopping server")
		s.App.Shutdown()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Fatal().Err(err).Msg("server failed")
	case <-stopChan:
	}
}
