package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceparentHeader  = "traceparent"
	TracestateHeader   = "tracestate"
	GrpcTraceBinHeader = "Grpc-Trace-Bin"
)

var (
	traceContext propagation.TraceContext
)

type CanonicalMapCarrier propagation.MapCarrier

func (c CanonicalMapCarrier) Get(key string) string {
	return propagation.MapCarrier(c).Get(http.CanonicalHeaderKey(key))
}

func (c CanonicalMapCarrier) Set(key, value string) {
	propagation.MapCarrier(c).Set(http.CanonicalHeaderKey(key), value)
}

func (c CanonicalMapCarrier) Keys() []string {
	return propagation.MapCarrier(c).Keys()
}

type GRPCHeaderCarrier propagation.HeaderCarrier

func (c GRPCHeaderCarrier) Get(key string) string {
	if key == TraceparentHeader {
		grpcTraceBin := propagation.HeaderCarrier(c).Get(GrpcTraceBinHeader)
		sc, _ := SpanContextFromBinary([]byte(grpcTraceBin))
		return SpanContextToW3CString(sc)
	}
	return propagation.HeaderCarrier(c).Get(key)
}

func (c GRPCHeaderCarrier) Set(key, value string) {
	if key == TraceparentHeader {
		sc, _ := SpanContextFromW3CString(value)
		grpcTraceBin := BinaryFromSpanContext(sc)
		propagation.HeaderCarrier(c).Set(GrpcTraceBinHeader, string(grpcTraceBin))
	} else {
		propagation.HeaderCarrier(c).Set(key, value)
	}
}

func (c GRPCHeaderCarrier) Keys() []string {
	return propagation.HeaderCarrier(c).Keys()
}

func InjectTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) {
	traceContext.Inject(ctx, carrier)
}

func WithTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return traceContext.Extract(ctx, carrier)
}

func TraceContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}
