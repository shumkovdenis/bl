package trace

import (
	"encoding/hex"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

const (
	supportedVersion = 0
	maxVersion       = 254
)

var emptySpanContext trace.SpanContext

// BinaryFromSpanContext returns the binary format representation of a SpanContext.
//
// If sc is the zero value, Binary returns nil.
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

// SpanContextFromBinary returns the SpanContext represented by b.
//
// If b has an unsupported version ID or contains no TraceID, SpanContextFromBinary returns with ok==false.
func SpanContextFromBinary(b []byte) (sc trace.SpanContext, ok bool) {
	var scConfig trace.SpanContextConfig
	if len(b) == 0 || b[0] != 0 {
		return trace.SpanContext{}, false
	}
	b = b[1:]
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

// SpanContextToW3CString returns the SpanContext string representation.
func SpanContextToW3CString(sc trace.SpanContext) string {
	traceID := sc.TraceID()
	spanID := sc.SpanID()
	traceFlags := sc.TraceFlags()
	return fmt.Sprintf("%x-%x-%x-%x",
		[]byte{supportedVersion},
		traceID[:],
		spanID[:],
		[]byte{byte(traceFlags)})
}

// SpanContextFromW3CString extracts a span context from given string which got earlier from SpanContextToW3CString format.
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
