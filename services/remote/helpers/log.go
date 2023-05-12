package helpers

import (
	"log"
)

func logBegin(msg string) {
	log.Println("------------------------------")
	log.Println(msg)
	log.Println("------------------------------")
}

func logTraceHeader(header HeaderGetter, key string) {
	log.Printf("header %s=%s", key, GetHeader(header, key))
}

func logTraceHeaders(header HeaderGetter) {
	logTraceHeader(header, traceParentHeader)
	logTraceHeader(header, traceStateHeader)
	logTraceHeader(header, grpcTraceBinHeader)
}
