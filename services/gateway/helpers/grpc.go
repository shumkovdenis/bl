package helpers

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
)

func NewInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}

func NewLoggerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				logBegin("client request " + req.Spec().Procedure)
			} else {
				logBegin("server request " + req.Spec().Procedure)
			}
			logTraceHeaders(req.Header())

			res, err := next(ctx, req)

			if req.Spec().IsClient {
				logBegin("client response " + req.Spec().Procedure)
			} else {
				logBegin("server response " + req.Spec().Procedure)
			}
			logTraceHeaders(res.Header())

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func NewTraceInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				SetTraceHeader(ctx, req.Header(), traceParentHeader)
				SetTraceHeader(ctx, req.Header(), grpcTraceBinHeader)
			} else {
				ctx = WithAllTraceHeader(ctx, req.Header())
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func NewAppInterceptor(appID string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				SetAppIdHeader(req.Header(), appID)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
