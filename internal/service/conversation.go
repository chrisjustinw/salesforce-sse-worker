package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"salesforce-sse-worker/configs"
	"salesforce-sse-worker/internal/library"
	"salesforce-sse-worker/internal/model"
	"salesforce-sse-worker/internal/repository"
	"salesforce-sse-worker/internal/request"
	"salesforce-sse-worker/internal/response"
	"salesforce-sse-worker/internal/service/outbound"
)

type (
	ConversationService interface {
		GenerateToken(ctx context.Context, req request.GenerateTokenRequest) (string, error)
		CreateConversationProducer(ctx context.Context, req request.CreateConversationRequest) (string, error)
		CreateConversationConsumer(ctx context.Context, req request.CreateConversationRequest, partition int) error
	}

	ConversationServiceImpl struct {
		kafkaConfig                   configs.KafkaConfig
		kafkaProducer                 library.KafkaProducer
		salesforceOutbound            outbound.SalesforceOutbound
		conversationMappingRepository repository.ConversationMappingRepository
	}
)

func NewConversationService(kafkaConfig configs.KafkaConfig, kafkaProducer library.KafkaProducer, salesforceOutbound outbound.SalesforceOutbound, conversationMappingRepository repository.ConversationMappingRepository) ConversationService {
	return &ConversationServiceImpl{
		kafkaConfig:                   kafkaConfig,
		kafkaProducer:                 kafkaProducer,
		salesforceOutbound:            salesforceOutbound,
		conversationMappingRepository: conversationMappingRepository,
	}
}

func (m *ConversationServiceImpl) GenerateToken(ctx context.Context, req request.GenerateTokenRequest) (string, error) {
	for partition := 0; partition <= m.kafkaConfig.PartitionCount; partition++ {
		resp, err := m.salesforceOutbound.GenerateToken(ctx, req)
		if err != nil {
			continue
		}

		var data response.GenerateTokenResponse
		if err := json.Unmarshal(resp, &data); err != nil {
			continue
		}

		if _, err = m.conversationMappingRepository.Upsert(ctx, model.ConversationMapping{
			Partition: partition,
			Token:     data.AccessToken,
		}); err != nil {
			continue
		}
	}

	return "Successfully generated tokens", nil
}

func (m *ConversationServiceImpl) CreateConversationProducer(ctx context.Context, req request.CreateConversationRequest) (string, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal req: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: m.kafkaConfig.Topics[0],
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := m.kafkaProducer.Produce(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("failed to produce Kafka message: %w", err)
	}

	slog.InfoContext(ctx, "Kafka message produced",
		slog.String("conversationId", req.ConversationId),
		slog.String("topic", msg.Topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)

	return "Message successfully queued", nil
}

func (m *ConversationServiceImpl) CreateConversationConsumer(ctx context.Context, req request.CreateConversationRequest, partition int) error {
	conversationMapping, err := m.conversationMappingRepository.FindOneByPartition(ctx, partition)
	if conversationMapping == nil || err != nil {
		return fmt.Errorf("failed to find token for partition %d: %w", partition, err)
	}

	_, err = m.salesforceOutbound.CreateConversation(ctx, conversationMapping.Token, req)
	if err != nil {
		return fmt.Errorf("failed to create conversation in Salesforce: %w", err)
	}

	return nil
}
