package connect

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/shumkovdenis/bl/trace"
)

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

func AddDaprAppID(appID string) connect.UnaryInterceptorFunc {
	return AddHeader("dapr-app-id", appID)
}

func AddTraceContext() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				trace.InjectTraceContext(ctx, req.Header())
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func AddBinaryTraceContext() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				carrier := ConnectHeaderCarrier(req.Header())
				trace.InjectBinaryTraceContext(ctx, carrier)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func WithClientOptions(interceptors ...connect.Interceptor) connect.ClientOption {
	return connect.WithClientOptions(
		connect.WithGRPC(),
		connect.WithInterceptors(interceptors...),
	)
}
