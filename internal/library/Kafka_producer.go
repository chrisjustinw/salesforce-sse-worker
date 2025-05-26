package library

import (
	"context"
	"github.com/IBM/sarama"
	"salesforce-sse-worker/configs"
)

type (
	KafkaProducer interface {
		Produce(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
	}

	KafkaProducerImpl struct {
		syncProducer sarama.SyncProducer
	}
)

func NewKafkaProducer(cfg configs.KafkaConfig, saramaCfg *sarama.Config) (KafkaProducer, error) {
	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &KafkaProducerImpl{syncProducer: syncProducer}, nil
}

func (p *KafkaProducerImpl) Produce(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return p.syncProducer.SendMessage(msg)
}
