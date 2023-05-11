package helpers

import "context"

const (
	traceParentHeader  = "traceparent"
	traceStateHeader   = "tracestate"
	grpcTraceBinHeader = "grpc-trace-bin"
)

type traceContextKey string

const (
	traceparentContextKey  = traceContextKey(traceParentHeader)
	tracestateContextKey   = traceContextKey(traceStateHeader)
	grpcTraceBinContextKey = traceContextKey(grpcTraceBinHeader)
)

func WithTraceparent(ctx context.Context, traceparent string) context.Context {
	return context.WithValue(ctx, traceparentContextKey, traceparent)
}

func WithTracestate(ctx context.Context, tracestate string) context.Context {
	return context.WithValue(ctx, tracestateContextKey, tracestate)
}

func WithGRPCTraceBin(ctx context.Context, grpcTraceBin string) context.Context {
	return context.WithValue(ctx, grpcTraceBinContextKey, grpcTraceBin)
}

func WithTrace(ctx context.Context, traceparent, tracestate, grpcTraceBin string) context.Context {
	ctx = WithTraceparent(ctx, traceparent)
	ctx = WithTracestate(ctx, tracestate)
	ctx = WithGRPCTraceBin(ctx, grpcTraceBin)
	return ctx
}

func ExtractTrace(ctx context.Context, key traceContextKey) string {
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}
	return ""
}

func ExtractTraceparent(ctx context.Context) string {
	return ExtractTrace(ctx, traceparentContextKey)
}

func ExtractTracestate(ctx context.Context) string {
	return ExtractTrace(ctx, tracestateContextKey)
}

func ExtractGRPCTraceBin(ctx context.Context) string {
	return ExtractTrace(ctx, grpcTraceBinContextKey)
}
