package client

import (
	"context"
	"fmt"
	"net/http"
	"salesforce-sse-worker/internal/request"
)

type (
	HTTPClient interface {
		Get(ctx context.Context, request request.HTTPRequest) (*http.Response, error)
		Post(ctx context.Context, request request.HTTPRequest) (*http.Response, error)
	}

	HTTPClientImpl struct {
		HTTPClient http.Client
	}
)

func NewHTTPClient() HTTPClient {
	return &HTTPClientImpl{HTTPClient: http.Client{}}
}

func (h *HTTPClientImpl) Get(ctx context.Context, request request.HTTPRequest) (*http.Response, error) {
	return h.do(ctx, http.MethodGet, request)
}

func (h *HTTPClientImpl) Post(ctx context.Context, request request.HTTPRequest) (*http.Response, error) {
	return h.do(ctx, http.MethodPost, request)
}

func (h *HTTPClientImpl) do(ctx context.Context, method string, request request.HTTPRequest) (*http.Response, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, method, request.Path, request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	for key, val := range request.Headers {
		httpRequest.Header.Add(key, val)
	}

	queries := httpRequest.URL.Query()
	for key, val := range request.Queries {
		queries.Add(key, val)
	}
	httpRequest.URL.RawQuery = queries.Encode()

	return h.HTTPClient.Do(httpRequest)
}
