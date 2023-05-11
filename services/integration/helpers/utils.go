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

func setNewGRPCTraceHeaderFromContext(header HeaderSetter, ctx context.Context) {
	value := ExtractGRPCTraceBin(ctx)
	log.Println("ExtractGRPCTraceBin", value)
	if value != "" {
		sc, ok := utils.SpanContextFromBinary([]byte(value))
		log.Println("SpanContextFromBinary", sc, ok)
		if ok {
			newsc := sc.WithSpanID(newSpanID())
			val := utils.BinaryFromSpanContext(newsc)
			log.Println("BinaryFromSpanContext", string(val))
			header.Set(grpcTraceBinHeader, string(val))
		}
	}
}
