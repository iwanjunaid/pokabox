package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConfig interface {
	GetConfigMap() *kafka.ConfigMap
}
