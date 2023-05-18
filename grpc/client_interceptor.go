package grpc

import (
	"context"
	"net/http"

	"github.com/shumkovdenis/bl/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AddMetadata(key, value string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, key, value)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func AddDaprAppID(appID string) grpc.UnaryClientInterceptor {
	return AddMetadata("dapr-app-id", appID)
}

func AddTraceContext() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var carrier http.Header
		trace.InjectTraceContext(ctx, carrier)
		ctx = metadata.AppendToOutgoingContext(ctx,
			trace.TraceparentHeader, carrier.Get(trace.TraceparentHeader))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func AddBinaryTraceContext() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		carrier := make(http.Header)
		trace.InjectBinaryTraceContext(ctx, carrier)
		ctx = metadata.AppendToOutgoingContext(ctx,
			trace.GrpcTraceBinHeader, carrier.Get(trace.GrpcTraceBinHeader))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func WithClientOptions(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}
