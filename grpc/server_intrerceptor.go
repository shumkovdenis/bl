package grpc

import (
	"context"

	"github.com/shumkovdenis/bl/logger"
	"github.com/shumkovdenis/bl/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func InjectTraceContext() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			carrier := MetadataCarrier(md)
			ctx = trace.ExtractBinaryTraceContext(ctx, carrier)
		}
		return handler(ctx, req)
	}
}

func InjectTraceContextLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = logger.WithTraceContext(ctx)
		return handler(ctx, req)
	}
}

func WithServerOptions() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		InjectTraceContext(),
		InjectTraceContextLogger(),
	)
}
