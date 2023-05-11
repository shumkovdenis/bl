package helpers

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

func SetAppIdHeader(headers HeaderSetter, appID string) {
	headers.Set("dapr-app-id", appID)
}
