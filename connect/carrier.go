package connect

import (
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/shumkovdenis/bl/trace"
)

type ConnectHeaderCarrier http.Header

func (c ConnectHeaderCarrier) Get(key string) string {
	value := http.Header(c).Get(key)
	if key == trace.GrpcTraceBinHeader {
		b, _ := connect.DecodeBinaryHeader(value)
		return string(b)
	}
	return value
}

func (c ConnectHeaderCarrier) Set(key, value string) {
	if key == trace.GrpcTraceBinHeader {
		value = connect.EncodeBinaryHeader([]byte(value))
	}
	http.Header(c).Set(key, value)
}
