package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"salesforce-sse-worker/internal/library"
	"salesforce-sse-worker/internal/repository"
	"salesforce-sse-worker/internal/request"
	"salesforce-sse-worker/internal/service"
	"salesforce-sse-worker/internal/service/outbound"
)

type (
	KafkaHandlerImpl struct {
		conversationService           service.ConversationService
		salesforceOutbound            outbound.SalesforceOutbound
		conversationMappingRepository repository.ConversationMappingRepository
	}
)

func NewKafkaHandler(conversationService service.ConversationService, salesforceOutbound outbound.SalesforceOutbound, conversationMappingRepository repository.ConversationMappingRepository) library.KafkaHandler {

	return &KafkaHandlerImpl{
		conversationService:           conversationService,
		salesforceOutbound:            salesforceOutbound,
		conversationMappingRepository: conversationMappingRepository,
	}
}

func (c *KafkaHandlerImpl) Setup(ctx context.Context, claims map[string][]int32) error {
	for _, partitions := range claims {
		for _, partition := range partitions {
			slog.InfoContext(ctx, "SSE Subscribed", slog.Any("partition", partition))
			c.subscribe(ctx, partition)
		}
	}

	return nil
}

func (c *KafkaHandlerImpl) Cleanup(ctx context.Context, claims map[string][]int32) error {
	for _, partitions := range claims {
		for _, partition := range partitions {
			slog.InfoContext(ctx, "SSE Revoked", slog.Any("partition", partition))
		}
	}

	return nil
}

func (c *KafkaHandlerImpl) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	slog.InfoContext(ctx, "Kafka message consumed",
		slog.String("topic", message.Topic),
		slog.Int("partition", int(message.Partition)),
		slog.Any("value", message.Value),
	)

	var req request.CreateConversationRequest
	if err := json.Unmarshal(message.Value, &req); err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}

	return c.conversationService.CreateConversationConsumer(ctx, req, int(message.Partition))
}

func (c *KafkaHandlerImpl) subscribe(ctx context.Context, partition int32) {
	go func() {
		conversationMapping, err := c.conversationMappingRepository.FindOneByPartition(ctx, int(partition))
		if err != nil {
			return
		}

		if conversationMapping == nil {
			return
		}

		if err := c.salesforceOutbound.Subscribe(ctx, conversationMapping.Token); err != nil {
			return
		}
	}()
}
