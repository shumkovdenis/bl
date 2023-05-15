package logger

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/shumkovdenis/bl/trace"
)

func WithTraceContextLogger(ctx context.Context) context.Context {
	sc := trace.TraceContextFromContext(ctx)

	logger := log.With().
		Str("trace_id", sc.TraceID().String()).
		Str("span_id", sc.SpanID().String()).
		Logger()

	return logger.WithContext(ctx)
}
