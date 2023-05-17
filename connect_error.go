package main

import (
	"errors"

	example "github.com/shumkovdenis/protobuf-schema/gen/example/v1"

	"github.com/bufbuild/connect-go"
)

func newCallError(err error) error {
	connectErr := connect.NewError(connect.CodeInternal, err)

	retryInfo := &example.RetryInfo{Count: 1}

	if detail, err := connect.NewErrorDetail(retryInfo); err == nil {
		connectErr.AddDetail(detail)
	}

	return connectErr
}

func extractCallError(err error) (*example.RetryInfo, bool) {
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		return nil, false
	}

	for _, detail := range connectErr.Details() {
		msg, valueErr := detail.Value()
		if valueErr != nil {
			// Usually, errors here mean that we don't have the schema for this
			// Protobuf message.
			continue
		}
		if retryInfo, ok := msg.(*example.RetryInfo); ok {
			return retryInfo, true
		}
	}

	return nil, false
}
