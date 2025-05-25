package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"salesforce-sse-worker/configs"
	"salesforce-sse-worker/internal/outbound/client"
	"salesforce-sse-worker/internal/request"
)

const (
	ssePath                       = "/eventrouter/v1/sse"
	createConversationPath        = "/iamessage/api/v2/conversation"
	generateContinuationTokenPath = "/iamessage/api/v2/authorization/continuation-access-token"
)

type (
	SalesforceOutbound interface {
		Subscribe(ctx context.Context, token string) error
		CreateConversation(ctx context.Context, token string, req request.CreateConversationRequest) ([]byte, error)
		GenerateContinuationToken(ctx context.Context, token string) ([]byte, error)
	}

	SalesforceOutboundImpl struct {
		httpClient       client.HTTPClient
		salesforceConfig configs.SalesforceConfig
	}
)

func NewSalesforceOutbound(httpClient client.HTTPClient, sseConfig configs.SalesforceConfig) SalesforceOutbound {
	return &SalesforceOutboundImpl{
		httpClient:       httpClient,
		salesforceConfig: sseConfig,
	}
}

func (s *SalesforceOutboundImpl) Subscribe(ctx context.Context, token string) error {
	url := s.salesforceConfig.Host + ssePath

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"X-Org-Id":      s.salesforceConfig.OrgId,
	}

	sseClient := client.NewSSEClient(url, headers)
	if err := sseClient.Start(ctx); err != nil {
		return err
	}

	return nil
}

func (s *SalesforceOutboundImpl) CreateConversation(ctx context.Context, token string, req request.CreateConversationRequest) ([]byte, error) {
	url := s.salesforceConfig.Host + createConversationPath

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post(ctx, request.HTTPRequest{
		Path:    url,
		Headers: headers,
		Body:    bytes.NewReader(payload),
	})

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *SalesforceOutboundImpl) GenerateContinuationToken(ctx context.Context, token string) ([]byte, error) {
	url := s.salesforceConfig.Host + generateContinuationTokenPath

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	resp, err := s.httpClient.Get(ctx, request.HTTPRequest{
		Path:    url,
		Headers: headers,
	})

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
