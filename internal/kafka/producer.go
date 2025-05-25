package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"salesforce-sse-worker/configs"
)

type (
	Producer interface {
		Produce(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
	}

	ProducerImpl struct {
		syncProducer sarama.SyncProducer
	}
)

func NewProducer(cfg configs.KafkaConfig, saramaCfg *sarama.Config) (Producer, error) {
	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &ProducerImpl{syncProducer: syncProducer}, nil
}

func (p *ProducerImpl) Produce(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return p.syncProducer.SendMessage(msg)
}
