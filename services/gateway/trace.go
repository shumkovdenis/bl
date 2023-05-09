package main

import "context"

const (
	TraceParentHeader  = "traceparent"
	TraceStateHeader   = "tracestate"
	GRPCTraceBinHeader = "grpc-trace-bin"
)

type traceContextKey string

const (
	traceparentContextKey  = traceContextKey(TraceParentHeader)
	tracestateContextKey   = traceContextKey(TraceStateHeader)
	grpcTraceBinContextKey = traceContextKey(GRPCTraceBinHeader)
)

func WithTraceparent(ctx context.Context, traceparent string) context.Context {
	return context.WithValue(ctx, traceparentContextKey, traceparent)
}

func ExtractTraceparent(ctx context.Context) string {
	if traceparent, ok := ctx.Value(traceparentContextKey).(string); ok {
		return traceparent
	}
	return ""
}

func WithTracestate(ctx context.Context, tracestate string) context.Context {
	return context.WithValue(ctx, tracestateContextKey, tracestate)
}

func ExtractTracestate(ctx context.Context) string {
	if tracestate, ok := ctx.Value(tracestateContextKey).(string); ok {
		return tracestate
	}
	return ""
}

func WithGRPCTraceBin(ctx context.Context, grpcTraceBin string) context.Context {
	return context.WithValue(ctx, grpcTraceBinContextKey, grpcTraceBin)
}

func ExtractGRPCTraceBin(ctx context.Context) string {
	if grpcTraceBin, ok := ctx.Value(grpcTraceBinContextKey).(string); ok {
		return grpcTraceBin
	}
	return ""
}

func WithTrace(ctx context.Context, traceparent, tracestate, grpcTraceBin string) context.Context {
	ctx = WithTraceparent(ctx, traceparent)
	ctx = WithTracestate(ctx, tracestate)
	ctx = WithGRPCTraceBin(ctx, grpcTraceBin)
	return ctx
}
