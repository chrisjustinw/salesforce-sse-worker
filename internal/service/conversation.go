package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"salesforce-sse-worker/internal/kafka"
	"salesforce-sse-worker/internal/model"
	"salesforce-sse-worker/internal/outbound"
	"salesforce-sse-worker/internal/repository"
	"salesforce-sse-worker/internal/request"
	"salesforce-sse-worker/internal/response"
)

type (
	ConversationService interface {
		PublishConversation(ctx context.Context, request request.CreateConversationRequest) (string, error)
		ConsumeConversation(ctx context.Context, partition int, message []byte) error
		GenerateContinuationToken(ctx context.Context) (string, error)
	}

	ConversationServiceImpl struct {
		kafkaProducer                 kafka.Producer
		salesforceOutbound            outbound.SalesforceOutbound
		conversationMappingRepository repository.ConversationMappingRepository
	}
)

func NewConversationService(kafkaProducer kafka.Producer, salesforceOutbound outbound.SalesforceOutbound, conversationMappingRepository repository.ConversationMappingRepository) ConversationService {
	return &ConversationServiceImpl{
		kafkaProducer:                 kafkaProducer,
		salesforceOutbound:            salesforceOutbound,
		conversationMappingRepository: conversationMappingRepository,
	}
}

func (m *ConversationServiceImpl) PublishConversation(ctx context.Context, request request.CreateConversationRequest) (string, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "conversation.create",
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := m.kafkaProducer.Produce(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("failed to produce Kafka message: %w", err)
	}

	slog.InfoContext(ctx, "Kafka message produced",
		slog.String("conversationId", request.ConversationId),
		slog.String("topic", msg.Topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)

	return "Message successfully queued", nil
}

func (m *ConversationServiceImpl) ConsumeConversation(ctx context.Context, partition int, message []byte) error {
	conversationMapping, err := m.conversationMappingRepository.FindOneByPartition(ctx, partition)
	if conversationMapping == nil || err != nil {
		return fmt.Errorf("failed to find token for partition %d: %w", partition, err)
	}

	slog.InfoContext(ctx, "Consuming conversation message",
		slog.Int("partition", partition),
		slog.String("token", conversationMapping.Token),
		slog.String("message", string(message)),
	)

	var req request.CreateConversationRequest
	if err := json.Unmarshal(message, &req); err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}

	_, err = m.salesforceOutbound.CreateConversation(ctx, conversationMapping.Token, req)
	if err != nil {
		return fmt.Errorf("failed to create conversation in Salesforce: %w", err)
	}

	return nil
}

func (m *ConversationServiceImpl) GenerateContinuationToken(ctx context.Context) (string, error) {
	conversationMappings, err := m.conversationMappingRepository.FindAll(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve conversation mappings: %w", err)
	}

	for _, mapping := range conversationMappings {
		resp, err := m.salesforceOutbound.GenerateContinuationToken(ctx, mapping.Token)
		if err != nil {
			continue
		}

		var data response.GenerateContinuationTokenResponse
		if err := json.Unmarshal(resp, &data); err != nil {
			continue
		}

		if _, err = m.conversationMappingRepository.Upsert(ctx, model.ConversationMapping{
			Partition: mapping.Partition,
			Token:     data.AccessToken,
		}); err != nil {
			continue
		}
	}

	return "Successfully generated continuation tokens", nil
}
