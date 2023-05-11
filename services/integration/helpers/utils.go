package helpers

import (
	"context"
	"log"
	"math/rand"

	"github.com/dapr/dapr/pkg/diagnostics/utils"
	"go.opentelemetry.io/otel/trace"
)

type HeaderGetter interface {
	Get(key string) string
}

type HeaderSetter interface {
	Set(key, value string)
}

func logHeader(key string, header HeaderGetter) {
	log.Println(key, header.Get(key))
}

func setTraceHeaderFromContext(key traceContextKey, header HeaderSetter, ctx context.Context) {
	value := ExtractTrace(ctx, key)
	if value != "" {
		header.Set(string(key), value)
	}
}

func newSpanID() trace.SpanID {
	randSource := rand.New(rand.NewSource(1))
	sid := trace.SpanID{}
	_, _ = randSource.Read(sid[:])
	return sid
}

func SpanContextFromBinary(b []byte) (sc trace.SpanContext, ok bool) {
	var scConfig trace.SpanContextConfig
	log.Println("SpanContextFromBinary", len(b), b[0])
	if len(b) == 0 || b[0] != 0 {
		return trace.SpanContext{}, false
	}
	b = b[1:]
	log.Println("SpanContextFromBinary", len(b), b[0])
	if len(b) >= 17 && b[0] == 0 {
		copy(scConfig.TraceID[:], b[1:17])
		b = b[17:]
	} else {
		return trace.SpanContext{}, false
	}
	if len(b) >= 9 && b[0] == 1 {
		copy(scConfig.SpanID[:], b[1:9])
		b = b[9:]
	}
	if len(b) >= 2 && b[0] == 2 {
		scConfig.TraceFlags = trace.TraceFlags(b[1])
	}
	sc = trace.NewSpanContext(scConfig)
	return sc, true
}

func setNewGRPCTraceHeaderFromContext(header HeaderSetter, ctx context.Context) {
	value := ExtractGRPCTraceBin(ctx)
	log.Println("ExtractGRPCTraceBin", value)
	if value != "" {
		sc, ok := SpanContextFromBinary([]byte(value))
		log.Println("SpanContextFromBinary", sc, ok)
		if ok {
			newsc := sc.WithSpanID(newSpanID())
			val := utils.BinaryFromSpanContext(newsc)
			log.Println("BinaryFromSpanContext", string(val))
			header.Set(grpcTraceBinHeader, string(val))
		}
	}
}
