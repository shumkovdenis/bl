package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shumkovdenis/bl/http/client"
)

type httpCallee struct {
	url    string
	client *http.Client
}

func NewHTTPCallee(cfg Config) *httpCallee {
	url := fmt.Sprintf("http://localhost:%d/call", cfg.Dapr.HTTPPort)

	client := client.NewClient(
		client.AddDaprAppIDHeader(cfg.Callee.ServiceName),
	)

	return &httpCallee{
		url:    url,
		client: client,
	}
}

func (c *httpCallee) Call(ctx context.Context, msg Message) (Message, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, nil)
	if err != nil {
		return Message{}, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return Message{}, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return Message{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var m Message
	if err := json.Unmarshal(b, &m); err != nil {
		return Message{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return m, nil
}
