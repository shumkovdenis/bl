package helpers

import (
	"context"
	"strings"
)

type traceContextKey string

func WithTraceHeader(ctx context.Context, header HeaderGetter, key string) context.Context {
	if value := strings.TrimSpace(GetHeader(header, key)); value != "" {
		ctx = context.WithValue(ctx, traceContextKey(key), value)
	}
	return ctx
}

func WithAllTraceHeader(ctx context.Context, header HeaderGetter) context.Context {
	ctx = WithTraceHeader(ctx, header, traceParentHeader)
	ctx = WithTraceHeader(ctx, header, traceStateHeader)
	ctx = WithTraceHeader(ctx, header, grpcTraceBinHeader)
	return ctx
}

func ExtractTraceHeader(ctx context.Context, key string) string {
	if value, ok := ctx.Value(traceContextKey(key)).(string); ok {
		return value
	}
	return ""
}

func SetTraceHeader(ctx context.Context, header HeaderSetter, key string) {
	value := ExtractTraceHeader(ctx, key)
	if value != "" {
		header.Set(key, value)
	}
}
