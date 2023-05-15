package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shumkovdenis/bl/logger"
	"github.com/shumkovdenis/bl/trace"
)

func InjectTraceContext() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = trace.WithTraceContextFromMap(ctx, c.GetReqHeaders())
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func InjectTraceContextLogger() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = logger.WithTraceContextLogger(ctx)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
