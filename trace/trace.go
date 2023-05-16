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
	GrpcTraceBinHeader = "grpc-trace-bin"
)

var (
	traceContext propagation.TraceContext
)

type canonicalMapCarrier propagation.MapCarrier

func (c canonicalMapCarrier) Get(key string) string {
	return propagation.MapCarrier(c).Get(http.CanonicalHeaderKey(key))
}

func (c canonicalMapCarrier) Set(key, value string) {
	propagation.MapCarrier(c).Set(http.CanonicalHeaderKey(key), value)
}

func (c canonicalMapCarrier) Keys() []string {
	return propagation.MapCarrier(c).Keys()
}

type grpcHeaderCarrier propagation.HeaderCarrier

func (c grpcHeaderCarrier) Get(key string) string {
	return propagation.HeaderCarrier(c).Get(http.CanonicalHeaderKey(key))
}

func (c grpcHeaderCarrier) Set(key, value string) {
	propagation.HeaderCarrier(c).Set(http.CanonicalHeaderKey(key), value)
}

func (c grpcHeaderCarrier) Keys() []string {
	return propagation.HeaderCarrier(c).Keys()
}

func WithTraceContextFromMap(ctx context.Context, headers map[string]string) context.Context {
	return traceContext.Extract(ctx, canonicalMapCarrier(headers))
}

func TraceContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}

func InjectTraceContext(ctx context.Context, header http.Header) {
	traceContext.Inject(ctx, propagation.HeaderCarrier(header))
}
