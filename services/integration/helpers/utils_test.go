package helpers

import (
	"encoding/base64"
	"testing"
)

const (
	traceparent  = "00-f1780472bc9d01296f455d1dbffa3351-e3efc05568f38f63-01"
	grpcTraceBin = "AAD/w6QpXa5Gi9HrIU/ZGb6hAf4DHe7VGxVzAgE"
)

func TestSpanContextFromBinary(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		b, _ := base64.StdEncoding.DecodeString(grpcTraceBin)
		t.Log(b)

		// sc, ok := SpanContextFromW3CString(traceparent)
		// t.Log(ok, sc.TraceID(), sc.SpanID(), sc.TraceFlags())

		// b := BinaryFromSpanContext(sc)
		// s := string(b)
		// t.Log(s)

		// sEnc := base64.StdEncoding.EncodeToString(b)
		// t.Log(sEnc)

		bsc, ok := SpanContextFromBinary(b)
		t.Log(ok, bsc.TraceID(), bsc.SpanID(), bsc.TraceFlags())

		// var nsc trace.SpanContext
		// t.Log(nsc.TraceID(), nsc.SpanID(), nsc.IsSampled())

		// spanContext := SpanContextFromBinary
		// sc, ok := spanContext([]byte("AAD/w6QpXa5Gi9HrIU/ZGb6hAf4DHe7VGxVzAgE"))
		// t.Log(sc, ok)
	})
}

// AAD/w6QpXa5Gi9HrIU/ZGb6hAf4DHe7VGxVzAgE
// AADxeARyvJ0BKW9FXR2/+jNRAePvwFVo849jAgE=
