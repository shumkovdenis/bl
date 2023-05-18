package grpc

import "google.golang.org/grpc/metadata"

type MetadataCarrier metadata.MD

func (c MetadataCarrier) Get(key string) string {
	values := metadata.MD(c).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c MetadataCarrier) Set(key, value string) {
	metadata.MD(c).Set(key, value)
}
