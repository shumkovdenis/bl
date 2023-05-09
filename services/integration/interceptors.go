package main

import (
	"context"
	"log"

	"github.com/bufbuild/connect-go"
)

func NewTraceLoggerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			log.Println("logger interceptor", TraceParentHeader, req.Header().Get(TraceParentHeader), req.Spec().IsClient)
			log.Println("logger interceptor", TraceStateHeader, req.Header().Get(TraceStateHeader), req.Spec().IsClient)
			log.Println("logger interceptor", GRPTraceBinHeader, req.Header().Get(GRPTraceBinHeader), req.Spec().IsClient)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
