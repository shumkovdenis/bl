package main

import "context"

var fakeMsg = Message{Content: "fake"}

type fakeCallee struct{}

func NewFakeCallee() *fakeCallee {
	return &fakeCallee{}
}

func (c *fakeCallee) Call(ctx context.Context, msg Message) (Message, error) {
	return fakeMsg, nil
}
