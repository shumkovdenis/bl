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

// var (
// 	traceContext propagation.TraceContext
// )

// type CanonicalMapCarrier propagation.MapCarrier

// func (c CanonicalMapCarrier) Get(key string) string {
// 	return propagation.MapCarrier(c).Get(http.CanonicalHeaderKey(key))
// }

// func (c CanonicalMapCarrier) Set(key, value string) {
// 	propagation.MapCarrier(c).Set(http.CanonicalHeaderKey(key), value)
// }

// func (c CanonicalMapCarrier) Keys() []string {
// 	return propagation.MapCarrier(c).Keys()
// }

// type GRPCHeaderCarrier propagation.HeaderCarrier

// func (c GRPCHeaderCarrier) Get(key string) string {
// 	if key == TraceparentHeader {
// 		grpcTraceBin := propagation.HeaderCarrier(c).Get(GrpcTraceBinHeader)
// 		sc, _ := SpanContextFromBinary([]byte(grpcTraceBin))
// 		return SpanContextToW3CString(sc)
// 	}
// 	return propagation.HeaderCarrier(c).Get(key)
// }

// func (c GRPCHeaderCarrier) Set(key, value string) {
// 	if key == TraceparentHeader {
// 		sc, _ := SpanContextFromW3CString(value)
// 		grpcTraceBin := BinaryFromSpanContext(sc)
// 		propagation.HeaderCarrier(c).Set(GrpcTraceBinHeader, string(grpcTraceBin))
// 	} else {
// 		propagation.HeaderCarrier(c).Set(key, value)
// 	}
// }

// func (c GRPCHeaderCarrier) Keys() []string {
// 	return propagation.HeaderCarrier(c).Keys()
// }

// func InjectTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) {
// 	traceContext.Inject(ctx, carrier)
// }

// func WithTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
// 	return traceContext.Extract(ctx, carrier)
// }

// func TraceContextFromContext(ctx context.Context) trace.SpanContext {
// 	return trace.SpanContextFromContext(ctx)
// }

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
}

func TraceContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}
