package library

import (
	"context"
	"github.com/IBM/sarama"
	"log/slog"
	"salesforce-sse-worker/configs"
	"sync"
)

type (
	KafkaHandler interface {
		Setup(ctx context.Context, claims map[string][]int32) error
		Cleanup(ctx context.Context, claims map[string][]int32) error
		Handle(ctx context.Context, message *sarama.ConsumerMessage) error
	}

	KafkaConsumerHandler interface {
		sarama.ConsumerGroupHandler
	}

	KafkaConsumerHandlerImpl struct {
		handler KafkaHandler
	}

	KafkaConsumer interface {
		ResetReady()
		WaitReady()
		Consume(ctx context.Context)
	}

	KafkaConsumerImpl struct {
		ready           chan bool
		topics          []string
		consumerGroup   sarama.ConsumerGroup
		consumerHandler sarama.ConsumerGroupHandler
	}
)

func NewKafkaConsumerHandler(handler KafkaHandler) KafkaConsumerHandler {
	return &KafkaConsumerHandlerImpl{handler: handler}
}

func (c *KafkaConsumerHandlerImpl) Setup(session sarama.ConsumerGroupSession) error {
	return c.handler.Setup(session.Context(), session.Claims())
}

func (c *KafkaConsumerHandlerImpl) Cleanup(session sarama.ConsumerGroupSession) error {
	return c.handler.Cleanup(session.Context(), session.Claims())
}

func (c *KafkaConsumerHandlerImpl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		sdc := map[string]string{"topic": message.Topic, "partition": string(message.Partition)}
		slog.InfoContext(context.Background(), "Message claimed", slog.Any("sdc", sdc))

		if err := c.handler.Handle(context.Background(), message); err != nil {
			slog.ErrorContext(session.Context(), "Failed to handle message", slog.Any("error", err))
		}

		session.MarkMessage(message, "")
	}

	return nil
}

func NewKafkaConsumer(cfg configs.KafkaConfig, saramaCfg *sarama.Config, handler KafkaHandler) (KafkaConsumer, error) {
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupName, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerImpl{
		ready:           make(chan bool),
		topics:          cfg.Topics,
		consumerGroup:   consumerGroup,
		consumerHandler: NewKafkaConsumerHandler(handler),
	}, nil
}

func (c *KafkaConsumerImpl) ResetReady() {
	c.ready = make(chan bool)
}

func (c *KafkaConsumerImpl) WaitReady() {
	<-c.ready
}

func (c *KafkaConsumerImpl) Consume(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			if ctx.Err() != nil {
				slog.InfoContext(ctx, "Context cancelled, stopping consumer")
				return
			}

			if err := c.consumerGroup.Consume(ctx, c.topics, c.consumerHandler); err != nil {
				slog.ErrorContext(ctx, "Error while consuming", slog.Any("error", err))
			}

			c.ResetReady()
		}
	}()

	c.WaitReady()
}
