package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shumkovdenis/bl/logger"
	"github.com/shumkovdenis/bl/trace"
)

func InjectTraceContext() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		carrier := trace.MapCarrier(c.GetReqHeaders())
		ctx = trace.ExtractTraceContext(ctx, carrier)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func InjectTraceContextLogger() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = logger.WithTraceContext(ctx)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
