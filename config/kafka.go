package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
)

type CommonKafkaConfig struct {
	configMap *kafka.ConfigMap
}

func NewCommonKafkaConfig(configMap *kafka.ConfigMap) config.KafkaConfig {
	config := &CommonKafkaConfig{
		configMap,
	}

	return config
}

func (c *CommonKafkaConfig) GetConfigMap() *kafka.ConfigMap {
	return c.configMap
}
