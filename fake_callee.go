package main

import "context"

type fakeCallee struct{}

func NewFakeCallee() *fakeCallee {
	return &fakeCallee{}
}

func (c fakeCallee) Call(ctx context.Context, msg Message) (Message, error) {
	return msg, nil
}
