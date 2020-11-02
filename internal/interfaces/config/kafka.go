package config

type KafkaConfig interface {
	GetBootstrapServers() string
}
