package main

import "context"

type Message struct {
	Content string `json:"content"`
}

type Callee interface {
	Call(ctx context.Context, msg Message) (Message, error)
}
