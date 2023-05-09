package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func NewTraceLoggerMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Println("logger middleware", TraceParentHeader, c.Get(TraceParentHeader))
		log.Println("logger middleware", TraceStateHeader, c.Get(TraceStateHeader))
		log.Println("logger middleware", GRPTraceBinHeader, c.Get(GRPTraceBinHeader))
		return c.Next()
	}
}
