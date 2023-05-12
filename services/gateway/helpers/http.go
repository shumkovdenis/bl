package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
)

type headerAdapter struct {
	m map[string]string
}

func (a headerAdapter) Get(key string) string {
	if v, ok := a.m[key]; ok {
		return v
	}
	return ""
}

func NewServerLoggerMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logBegin("server request " + c.Path())
		logTraceHeaders(&headerAdapter{m: c.GetReqHeaders()})

		err := c.Next()

		logBegin("server response " + c.Path())
		logTraceHeaders(&headerAdapter{m: c.GetRespHeaders()})

		return err
	}
}

func NewServerTraceMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := WithAllTraceHeader(c.UserContext(), &headerAdapter{m: c.GetReqHeaders()})

		c.SetUserContext(ctx)

		return c.Next()
	}
}

func NewClientLoggerMiddleware() func(rt req.RoundTripper) req.RoundTripFunc {
	return func(rt req.RoundTripper) req.RoundTripFunc {
		return func(req *req.Request) (*req.Response, error) {
			logBegin("client request " + req.URL.Path)
			logTraceHeaders(req.Headers)

			res, err := rt.RoundTrip(req)

			logBegin("client response " + req.URL.Path)
			logTraceHeaders(res.Header)

			return res, err
		}
	}
}

func NewClientTraceMiddleware() func(rt req.RoundTripper) req.RoundTripFunc {
	return func(rt req.RoundTripper) req.RoundTripFunc {
		return func(req *req.Request) (*req.Response, error) {
			SetTraceHeader(req.Context(), req.Headers, traceParentHeader)
			SetTraceHeader(req.Context(), req.Headers, traceStateHeader)

			return rt.RoundTrip(req)
		}
	}
}

func NewClientAppMiddleware(appID string) func(rt req.RoundTripper) req.RoundTripFunc {
	return func(rt req.RoundTripper) req.RoundTripFunc {
		return func(req *req.Request) (*req.Response, error) {
			SetAppIdHeader(req.Headers, appID)

			return rt.RoundTrip(req)
		}
	}
}
