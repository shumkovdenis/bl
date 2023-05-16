package connect

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/shumkovdenis/bl/logger"
	"github.com/shumkovdenis/bl/trace"
)

type ConnectHeaderCarrier http.Header

func (c ConnectHeaderCarrier) Get(key string) string {
	value := http.Header(c).Get(key)
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
	http.Header(c).Set(key, value)
}

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
				trace.InjectTraceContext(ctx, req.Header())
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func AddBinaryTraceContextHeader() connect.UnaryInterceptorFunc {
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
