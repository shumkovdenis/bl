package trace

import "go.opentelemetry.io/otel/trace"

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
