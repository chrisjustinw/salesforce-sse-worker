package configs

import (
	"github.com/IBM/sarama"
	"github.com/kelseyhightower/envconfig"
)

type KafkaConfig struct {
	Brokers   []string `envconfig:"BROKERS"`
	Topics    []string `envconfig:"TOPICS"`
	GroupName string   `envconfig:"GROUP_NAME"`
}

func NewKafkaConfig(e EnvFileRead) (KafkaConfig, error) {
	var cfg KafkaConfig
	if err := envconfig.Process("KAFKA", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func NewSaramaConfig(e EnvFileRead) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	return saramaConfig
}
