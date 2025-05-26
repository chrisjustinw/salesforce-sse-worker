package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"salesforce-sse-worker/configs"
	"salesforce-sse-worker/internal/library"
	"salesforce-sse-worker/internal/request"
)

const (
	ssePath                = "/eventrouter/v1/sse"
	createConversationPath = "/iamessage/api/v2/conversation"
	generateTokenPath      = "/iamessage/api/v2/authorization/unauthenticated/access-token"
)

type (
	SalesforceOutbound interface {
		GenerateToken(ctx context.Context, req request.GenerateTokenRequest) ([]byte, error)
		CreateConversation(ctx context.Context, token string, req request.CreateConversationRequest) ([]byte, error)
		Subscribe(ctx context.Context, token string) error
	}

	SalesforceOutboundImpl struct {
		httpClient       library.HTTPClient
		salesforceConfig configs.SalesforceConfig
	}
)

func NewSalesforceOutbound(httpClient library.HTTPClient, sseConfig configs.SalesforceConfig) SalesforceOutbound {
	return &SalesforceOutboundImpl{
		httpClient:       httpClient,
		salesforceConfig: sseConfig,
	}
}

func (s *SalesforceOutboundImpl) GenerateToken(ctx context.Context, req request.GenerateTokenRequest) ([]byte, error) {
	url := s.salesforceConfig.Host + generateTokenPath

	headers := map[string]string{
		"Content-Type": "application/json",
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

func (s *SalesforceOutboundImpl) Subscribe(ctx context.Context, token string) error {
	url := s.salesforceConfig.Host + ssePath

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"X-Org-Id":      s.salesforceConfig.OrgId,
	}

	sseClient := library.NewSSEClient(url, headers)
	if err := sseClient.Start(ctx); err != nil {
		return err
	}

	return nil
}
