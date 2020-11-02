package config

import (
	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
)

type CommonKafkaConfig struct {
	BootstrapServers string
}

func NewCommonKafkaConfig(bootstrapServers string) config.KafkaConfig {
	config := &CommonKafkaConfig{
		bootstrapServers,
	}

	return config
}

func (c *CommonKafkaConfig) GetBootstrapServers() string {
	return c.BootstrapServers
}
