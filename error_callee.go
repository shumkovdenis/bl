package main

import (
	"context"
	"errors"
)

type errorCallee struct{}

func NewErrorCallee() *errorCallee {
	return &errorCallee{}
}

func (c errorCallee) Call(ctx context.Context, msg Message) (Message, error) {
	return Message{}, errors.New("error in callee")
}
