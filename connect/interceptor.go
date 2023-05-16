package connect

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/shumkovdenis/bl/logger"
	"github.com/shumkovdenis/bl/trace"
	"go.opentelemetry.io/otel/propagation"
)

type ConnectHeaderCarrier propagation.HeaderCarrier

func (c ConnectHeaderCarrier) Get(key string) string {
	value := propagation.HeaderCarrier(c).Get(key)
	if key == trace.GrpcTraceBinHeader {
		b, _ := connect.DecodeBinaryHeader(value)
		return string(b)
	}
	return value
}

func (c ConnectHeaderCarrier) Set(key, value string) {
	if key == trace.GrpcTraceBinHeader {
		value = connect.EncodeBinaryHeader([]byte(value))
	}
	propagation.HeaderCarrier(c).Set(key, value)
}

func (c ConnectHeaderCarrier) Keys() []string {
	return propagation.HeaderCarrier(c).Keys()
}

func InjectTraceContext() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if !req.Spec().IsClient {
				carrier := trace.GRPCHeaderCarrier(ConnectHeaderCarrier(req.Header()))
				ctx = trace.WithTraceContext(ctx, carrier)
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

func AddTraceContextHeader() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				header := req.Header()
				carrier := trace.GRPCHeaderCarrier(
					ConnectHeaderCarrier(header))
				trace.InjectTraceContext(ctx, carrier)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
