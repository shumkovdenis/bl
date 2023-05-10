package helpers

import (
	"context"
	"log"
	"net/http"
)

func setHeaderFromContext(key traceContextKey, header http.Header, ctx context.Context) {
	value := ExtractTrace(ctx, key)
	log.Println("set header from context", key, value)
	if value != "" {
		header.Set(string(key), value)
	}
}

func logHeader(key string, header http.Header) {
	log.Println(key, header.Get(key))
}
