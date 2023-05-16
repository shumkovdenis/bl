package client

import (
	"net/http"

	"github.com/shumkovdenis/bl/trace"
)

// https://jonfriesen.ca/articles/go-http-client-middleware

// internalRoundTripper is a holder function to make the process of
// creating middleware a bit easier without requiring the consumer to
// implement the RoundTripper interface.
type internalRoundTripper func(*http.Request) (*http.Response, error)

func (rt internalRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// Middleware is our middleware creation functionality.
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain is a handy function to wrap a base RoundTripper (optional)
// with the middlewares.
func Chain(rt http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	for _, m := range middlewares {
		rt = m(rt)
	}

	return rt
}

// AddHeader adds a header to the request.
func AddHeader(key, value string) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRoundTripper(func(req *http.Request) (*http.Response, error) {
			header := req.Header
			if header == nil {
				header = make(http.Header)
			}

			header.Set(key, value)

			return rt.RoundTrip(req)
		})
	}
}

func AddDaprAppIDHeader(appID string) Middleware {
	return AddHeader("dapr-app-id", appID)
}

func AddTraceContextHeader() Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRoundTripper(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			header := req.Header
			if header == nil {
				header = make(http.Header)
			}

			trace.InjectTraceContext(ctx, header)

			return rt.RoundTrip(req)
		})
	}
}
