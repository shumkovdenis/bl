package helpers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type headerAdapter struct {
	c *fiber.Ctx
}

func (a *headerAdapter) Get(key string) string {
	return a.c.Get(key)
}

func NewLoggerMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		adapter := headerAdapter{c: c}
		log.Println("------------------------------")
		log.Println("logger middleware", c.Path())
		log.Println("------------------------------")
		logHeader(traceParentHeader, &adapter)
		logHeader(traceStateHeader, &adapter)
		logHeader(grpcTraceBinHeader, &adapter)

		return c.Next()
	}
}

func NewTraceMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := WithTrace(c.UserContext(),
			c.Get(traceParentHeader),
			c.Get(traceStateHeader),
			c.Get(grpcTraceBinHeader),
		)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
