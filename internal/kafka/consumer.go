package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log/slog"
	"salesforce-sse-worker/configs"
	"sync"
)

type (
	Handler interface {
		Setup(ctx context.Context, claims map[string][]int32) error
		Cleanup(ctx context.Context, claims map[string][]int32) error
		Handle(ctx context.Context, message *sarama.ConsumerMessage) error
	}

	Consumer interface {
		ResetReady()
		WaitReady()
		Consume(ctx context.Context)
	}

	ConsumerImpl struct {
		ready           chan bool
		topics          []string
		consumerGroup   sarama.ConsumerGroup
		consumerHandler sarama.ConsumerGroupHandler
	}

	ConsumerHandler interface {
		sarama.ConsumerGroupHandler
	}

	ConsumerHandlerImpl struct {
		handler Handler
	}
)

func NewConsumer(cfg configs.KafkaConfig, saramaCfg *sarama.Config, handler Handler) (Consumer, error) {
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupName, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &ConsumerImpl{
		ready:           make(chan bool),
		topics:          cfg.Topics,
		consumerGroup:   consumerGroup,
		consumerHandler: NewConsumerHandler(handler),
	}, nil
}

func (c *ConsumerImpl) ResetReady() {
	c.ready = make(chan bool)
}

func (c *ConsumerImpl) WaitReady() {
	<-c.ready
}

func (c *ConsumerImpl) Consume(ctx context.Context) {
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

func NewConsumerHandler(handler Handler) ConsumerHandler {
	return &ConsumerHandlerImpl{handler: handler}
}

func (c *ConsumerHandlerImpl) Setup(session sarama.ConsumerGroupSession) error {
	return c.handler.Setup(session.Context(), session.Claims())
}

func (c *ConsumerHandlerImpl) Cleanup(session sarama.ConsumerGroupSession) error {
	return c.handler.Cleanup(session.Context(), session.Claims())
}

func (c *ConsumerHandlerImpl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
