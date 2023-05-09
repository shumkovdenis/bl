package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CopyTraceHeaders(ctx *fiber.Ctx, header http.Header) {
	header.Set("traceparent", ctx.Get("traceparent"))
	header.Set("tracestate", ctx.Get("tracestate"))
	header.Set("grpc-trace-bin", ctx.Get("grpc-trace-bin"))
}
