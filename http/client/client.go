package client

import "net/http"

func NewClient(middlewares ...Middleware) *http.Client {
	var m []Middleware
	m = append(m,
		AddTraceContextHeader(),
		AddHeader("Content-Type", "application/json"))
	m = append(m, middlewares...)

	return &http.Client{
		Transport: Chain(nil, m...),
	}
}
