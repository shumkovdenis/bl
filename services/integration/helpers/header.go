package helpers

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

func SetAppIdHeader(header HeaderSetter, appID string) {
	header.Set("dapr-app-id", appID)
}
