package helpers

import (
	"context"
	"log"

	"github.com/bufbuild/connect-go"
)

func NewLoggerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			logHeader := func(header string) {
				log.Println(header, req.Header().Get(header))
			}

			log.Println("logger interceptor where client is", req.Spec().IsClient)
			logHeader(TraceParentHeader)
			logHeader(TraceStateHeader)
			logHeader(GRPCTraceBinHeader)

			return next(ctx, req)
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
			// setHeader := func(key traceContextKey) {
			// 	value := ExtractTrace(ctx, key)
			// 	if value != "" {
			// 		req.Header().Set(string(key), value)
			// 	}
			// }

			if req.Spec().IsClient {
				// setHeader(traceparentContextKey)
				// setHeader(tracestateContextKey)
				// setHeader(grpcTraceBinContextKey)
			} else {
				ctx = WithTrace(ctx,
					req.Header().Get(TraceParentHeader),
					req.Header().Get(TraceStateHeader),
					req.Header().Get(GRPCTraceBinHeader),
				)
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
				req.Header().Set("dapr-app-id", appID)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
