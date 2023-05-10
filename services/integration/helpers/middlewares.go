package helpers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func NewLoggerMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logHeader := func(header string) {
			log.Println(header, c.Get(header))
		}

		log.Println("logger middleware", c.Path())
		logHeader(TraceParentHeader)
		logHeader(TraceStateHeader)
		logHeader(GRPCTraceBinHeader)

		return c.Next()
	}
}

func NewTraceMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := WithTrace(c.UserContext(),
			c.Get(TraceParentHeader),
			c.Get(TraceStateHeader),
			c.Get(GRPCTraceBinHeader),
		)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
