package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	TraceparentHeader  = "Traceparent"
	TracestateHeader   = "Tracestate"
	GrpcTraceBinHeader = "Grpc-Trace-Bin"
)

func ExtractTraceContext(ctx context.Context, carrier Carrier) context.Context {
	traceparent := carrier.Get(TraceparentHeader)
	sc, _ := SpanContextFromW3CString(traceparent)
	if !sc.IsValid() {
		return ctx
	}

	tracestate := carrier.Get(TracestateHeader)
	ts := TraceStateFromW3CString(tracestate)
	sc = sc.WithTraceState(ts)

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func InjectTraceContext(ctx context.Context, carrier Carrier) {
	sc := trace.SpanContextFromContext(ctx)

	traceparent := SpanContextToW3CString(sc)
	carrier.Set(TraceparentHeader, traceparent)

	tracestate := TraceStateToW3CString(sc)
	carrier.Set(TracestateHeader, tracestate)
}

func ExtractBinaryTraceContext(ctx context.Context, carrier Carrier) context.Context {
	grpcTraceBin := carrier.Get(GrpcTraceBinHeader)
	sc, _ := SpanContextFromBinary([]byte(grpcTraceBin))
	if !sc.IsValid() {
		return ctx
	}

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func InjectBinaryTraceContext(ctx context.Context, carrier Carrier) {
	sc := trace.SpanContextFromContext(ctx)
	grpcTraceBin := BinaryFromSpanContext(sc)
	carrier.Set(GrpcTraceBinHeader, string(grpcTraceBin))
	InjectTraceContext(ctx, carrier)
}

func TraceContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}
