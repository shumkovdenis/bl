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

// TraceStateFromW3CString extracts a span tracestate from given string which got earlier from TraceStateFromW3CString format.
func TraceStateFromW3CString(h string) trace.TraceState {
	if h == "" {
		return trace.TraceState{}
	}

	ts, err := trace.ParseTraceState(h)
	if err != nil {
		return trace.TraceState{}
	}

	return ts
}

// TraceStateToW3CString extracts the TraceState from given SpanContext and returns its string representation.
func TraceStateToW3CString(sc trace.SpanContext) string {
	return sc.TraceState().String()
}
