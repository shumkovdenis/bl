package connect

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/shumkovdenis/bl/logger"
	"github.com/shumkovdenis/bl/trace"
)

func InjectTraceContext() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if !req.Spec().IsClient {
				carrier := ConnectHeaderCarrier(req.Header())
				ctx = trace.ExtractBinaryTraceContext(ctx, carrier)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func InjectTraceContextLogger() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if !req.Spec().IsClient {
				ctx = logger.WithTraceContext(ctx)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func WithServerOptions(interceptors ...connect.Interceptor) connect.HandlerOption {
	return connect.WithHandlerOptions(
		connect.WithInterceptors(
			InjectTraceContext(),
			InjectTraceContextLogger(),
		),
		connect.WithInterceptors(interceptors...),
	)
}
