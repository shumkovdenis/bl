package helpers

import (
	"context"
	"log"
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

func setHeaderFromContext(key traceContextKey, header HeaderSetter, ctx context.Context) {
	value := ExtractTrace(ctx, key)
	if value != "" {
		header.Set(string(key), value)
	}
}
