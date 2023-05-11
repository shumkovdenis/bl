package helpers

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"log"
	"math/rand"
	"strings"

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

const maxVersion = 254

func SpanContextFromW3CString(h string) (sc trace.SpanContext, ok bool) {
	if h == "" {
		return trace.SpanContext{}, false
	}
	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return trace.SpanContext{}, false
	}

	if len(sections[0]) != 2 {
		return trace.SpanContext{}, false
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return trace.SpanContext{}, false
	}
	version := int(ver[0])
	if version > maxVersion {
		return trace.SpanContext{}, false
	}

	if version == 0 && len(sections) != 4 {
		return trace.SpanContext{}, false
	}

	if len(sections[1]) != 32 {
		return trace.SpanContext{}, false
	}
	tid, err := trace.TraceIDFromHex(sections[1])
	if err != nil {
		return trace.SpanContext{}, false
	}
	sc = sc.WithTraceID(tid)

	if len(sections[2]) != 16 {
		return trace.SpanContext{}, false
	}
	sid, err := trace.SpanIDFromHex(sections[2])
	if err != nil {
		return trace.SpanContext{}, false
	}
	sc = sc.WithSpanID(sid)

	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 {
		return trace.SpanContext{}, false
	}
	sc = sc.WithTraceFlags(trace.TraceFlags(opts[0]))

	// Don't allow all zero trace or span ID.
	if sc.TraceID() == [16]byte{} || sc.SpanID() == [8]byte{} {
		return trace.SpanContext{}, false
	}

	return sc, true
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

var emptySpanContext trace.SpanContext

func BinaryFromSpanContext(sc trace.SpanContext) []byte {
	traceID := sc.TraceID()
	spanID := sc.SpanID()
	traceFlags := sc.TraceFlags()
	if sc.Equal(emptySpanContext) {
		return nil
	}
	var b [29]byte
	copy(b[2:18], traceID[:])
	b[18] = 1
	copy(b[19:27], spanID[:])
	b[27] = 2
	b[28] = uint8(traceFlags)
	return b[:]
}

func setNewGRPCTraceHeaderFromContext(header HeaderSetter, ctx context.Context) {
	value := ExtractGRPCTraceBin(ctx)
	log.Println("ExtractGRPCTraceBin", value)
	if value != "" {
		b, _ := base64.StdEncoding.DecodeString(value)
		sc, ok := SpanContextFromBinary(b)
		log.Println("SpanContextFromBinary", sc, ok)
		if ok {
			newsc := sc.WithSpanID(newSpanID())
			val := BinaryFromSpanContext(newsc)
			s := base64.StdEncoding.EncodeToString(val)
			log.Println("BinaryFromSpanContext", s)
			header.Set(grpcTraceBinHeader, string(val))
		}
	}
}
