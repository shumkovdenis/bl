package main

import "context"

type Message struct {
	Content string `json:"content"`
}

type Caller interface {
	Call(ctx context.Context, msg Message) (Message, error)
}

type fakeCaller struct{}

func NewFakeCaller() *fakeCaller {
	return &fakeCaller{}
}

var fakeMsg = Message{Content: "fake"}

func (c *fakeCaller) Call(ctx context.Context, msg Message) (Message, error) {
	return fakeMsg, nil
}
