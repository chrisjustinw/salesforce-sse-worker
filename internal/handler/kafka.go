package handler

import (
	"context"
	"github.com/IBM/sarama"
	"log/slog"
	"salesforce-sse-worker/internal/kafka"
	"salesforce-sse-worker/internal/outbound"
	"salesforce-sse-worker/internal/repository"
	"salesforce-sse-worker/internal/service"
)

type (
	KafkaHandlerImpl struct {
		conversationService           service.ConversationService
		salesforceOutbound            outbound.SalesforceOutbound
		conversationMappingRepository repository.ConversationMappingRepository
	}
)

func NewKafkaHandler(conversationService service.ConversationService, salesforceOutbound outbound.SalesforceOutbound, conversationMappingRepository repository.ConversationMappingRepository) kafka.Handler {

	return &KafkaHandlerImpl{
		conversationService:           conversationService,
		salesforceOutbound:            salesforceOutbound,
		conversationMappingRepository: conversationMappingRepository,
	}
}

func (c *KafkaHandlerImpl) Setup(ctx context.Context, claims map[string][]int32) error {
	for _, partitions := range claims {
		for _, partition := range partitions {
			slog.InfoContext(ctx, "Subscribed SSE", slog.Any("partition", partition))
			c.subscribe(ctx, partition)
		}
	}

	return nil
}

func (c *KafkaHandlerImpl) Cleanup(ctx context.Context, claims map[string][]int32) error {
	for _, partitions := range claims {
		for _, partition := range partitions {
			slog.InfoContext(ctx, "Revoked SSE", slog.Any("partition", partition))
		}
	}

	return nil
}

func (c *KafkaHandlerImpl) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	return c.conversationService.ConsumeConversation(ctx, int(message.Partition), message.Value)
}

func (c *KafkaHandlerImpl) subscribe(ctx context.Context, partition int32) {
	go func() {
		conversationMapping, err := c.conversationMappingRepository.FindOneByPartition(ctx, int(partition))
		if conversationMapping == nil || err != nil {
			return
		}

		if err := c.salesforceOutbound.Subscribe(ctx, conversationMapping.Token); err != nil {
			return
		}
	}()
}
