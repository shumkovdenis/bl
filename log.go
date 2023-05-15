package main

import (
	"log"
)

func logTraceHeader(header HeaderGetter, key string) {
	value := GetHeader(header, key)
	if value != "" {
		log.Printf("header %s=%s", key, value)
	}
}

func logTraceHeaders(msg string, header HeaderGetter) {
	log.Println(">>>", msg)
	logTraceHeader(header, traceParentHeader)
	logTraceHeader(header, traceStateHeader)
	logTraceHeader(header, grpcTraceBinHeader)
}
