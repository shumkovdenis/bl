package connect

import (
	"context"
	"log"

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
				log.Println(req.Header())
				log.Println(req.Header().Get(trace.TraceparentHeader))
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
				ctx = logger.WithTraceContextLogger(ctx)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func AddHeader(key, value string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				req.Header().Set(key, value)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func AddDaprAppIDHeader(appID string) connect.UnaryInterceptorFunc {
	return AddHeader("dapr-app-id", appID)
}
