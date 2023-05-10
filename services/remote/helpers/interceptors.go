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
			log.Println("------------------------------")
			if req.Spec().IsClient {
				log.Println("client logger interceptor")
			} else {
				log.Println("server logger interceptor")
			}
			log.Println("------------------------------")
			logHeader(traceParentHeader, req.Header())
			logHeader(traceStateHeader, req.Header())
			logHeader(grpcTraceBinHeader, req.Header())
			log.Println("------------------------------")

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
			if req.Spec().IsClient {
				setHeaderFromContext(traceparentContextKey, req.Header(), ctx)
			} else {
				ctx = WithTrace(ctx,
					req.Header().Get(traceParentHeader),
					req.Header().Get(traceStateHeader),
					req.Header().Get(grpcTraceBinHeader),
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
