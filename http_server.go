package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type headerAdapter struct {
	headers map[string]string
}

func newHeaderAdapter(headers map[string]string) *headerAdapter {
	return &headerAdapter{headers: headers}
}

func (a headerAdapter) Get(key string) string {
	if v, ok := a.headers[key]; ok {
		return v
	}
	return ""
}

func serverLoggerMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logTraceHeaders("server request "+c.Path(),
			newHeaderAdapter(c.GetReqHeaders()))

		err := c.Next()

		logTraceHeaders("server response "+c.Path(),
			newHeaderAdapter(c.GetRespHeaders()))

		return err
	}
}

func serverTraceMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		ctx = WithAllTraceHeader(ctx, newHeaderAdapter(c.GetReqHeaders()))

		c.SetUserContext(ctx)

		return c.Next()
	}
}

type httpServer struct {
	caller Caller
}

func NewHTTPServer(cfg Config, caller Caller) error {
	server := httpServer{caller: caller}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(
		serverTraceMiddleware(),
		serverLoggerMiddleware(),
	)
	app.Post("/call", server.Handler)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (s *httpServer) Handler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var msg Message

	if err := c.BodyParser(&msg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	msg, err := s.caller.Call(ctx, msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(&msg)
}
