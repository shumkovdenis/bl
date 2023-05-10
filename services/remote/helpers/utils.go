package helpers

import (
	"context"
	"log"
	"net/http"
)

func setHeaderFromContext(key traceContextKey, header http.Header, ctx context.Context) {
	value := ExtractTrace(ctx, key)
	if value != "" {
		header.Set(string(key), value)
	}
}

func logHeader(key string, header http.Header) {
	log.Println(header, header.Get(key))
}
