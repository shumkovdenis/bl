package main

import "net/http"

const (
	traceParentHeader  = "traceparent"
	traceStateHeader   = "tracestate"
	grpcTraceBinHeader = "grpc-trace-bin"
)

type HeaderGetter interface {
	Get(key string) string
}

type HeaderSetter interface {
	Set(key, value string)
}

func GetHeader(header HeaderGetter, key string) string {
	return header.Get(http.CanonicalHeaderKey(key))
}

func SetHeader(header HeaderSetter, key, value string) {
	header.Set(http.CanonicalHeaderKey(key), value)
}

func SetAppIDHeader(header HeaderSetter, appID string) {
	SetHeader(header, "dapr-app-id", appID)
}
