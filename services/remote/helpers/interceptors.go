package helpers

import (
	"context"
	"log"

	"github.com/bufbuild/connect-go"
)

func modeName(isClient bool) string {
	if isClient {
		return "client"
	}
	return "server"
}

func NewLoggerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {

			log.Println("--------------------")
			log.Println("logger interceptor", modeName(req.Spec().IsClient))
			log.Println("--------------------")
			log.Println("request headers:")
			logHeader(traceParentHeader, req.Header())
			logHeader(traceStateHeader, req.Header())
			logHeader(grpcTraceBinHeader, req.Header())

			res, err := next(ctx, req)

			log.Println("--------------------")
			log.Println("response headers:")
			logHeader(traceParentHeader, res.Header())
			logHeader(traceStateHeader, res.Header())
			logHeader(grpcTraceBinHeader, res.Header())
			log.Println("--------------------")

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
				// setHeader(traceparentContextKey)
				// setHeader(tracestateContextKey)
				// setHeader(grpcTraceBinContextKey)
			} else {
				ctx = WithTrace(ctx,
					req.Header().Get(traceParentHeader),
					req.Header().Get(traceStateHeader),
					req.Header().Get(grpcTraceBinHeader),
				)
			}

			res, err := next(ctx, req)

			if req.Spec().IsClient {
				// SetHeaderFromContext(traceParentHeader, req.Header(), ctx)
				// SetHeaderFromContext(traceStateHeader, req.Header(), ctx)
				// setHeaderFromContext(grpcTraceBinContextKey, res.Header(), ctx)
			} else {
				setHeaderFromContext(grpcTraceBinContextKey, res.Header(), ctx)
			}

			return res, err
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
