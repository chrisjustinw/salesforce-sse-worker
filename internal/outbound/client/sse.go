package client

import (
	"context"
	"github.com/r3labs/sse/v2"
	"log/slog"
)

type (
	SSEClient interface {
		Start(ctx context.Context) error
	}

	SSEClientImpl struct {
		url     string
		headers map[string]string
		client  *sse.Client
	}
)

func NewSSEClient(url string, headers map[string]string) SSEClient {
	client := sse.NewClient(url)
	for k, v := range headers {
		client.Headers[k] = v
	}
	return &SSEClientImpl{
		url:     url,
		headers: headers,
		client:  client,
	}
}

func (s *SSEClientImpl) Start(ctx context.Context) error {
	err := s.client.SubscribeWithContext(ctx, "", func(msg *sse.Event) {
		slog.InfoContext(ctx, "Received message", "url", s.url, "data", msg.Data)
	})

	if err != nil {
		slog.ErrorContext(ctx, "Failed to subscribe to SSE", "url", s.url, "error", err)
		return err
	}

	return nil
}
