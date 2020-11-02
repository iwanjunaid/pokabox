package config

import (
	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
)

type CommonOutboxConfig struct {
	GroupID                uuid.UUID
	OutboxTableName        string
	PickPollInterval       int
	PickZombiePollInterval int
	ZombieInterval         int
	TotalMessagePerPick    int
}

func NewCommonOutboxConfig(groupID uuid.UUID, outboxTableName string,
	pickPollInterval int, pickZombiePollInterval int,
	zombieInterval int, totalMessagePerPick int) config.OutboxConfig {
	config := &CommonOutboxConfig{
		groupID,
		outboxTableName,
		pickPollInterval,
		pickZombiePollInterval,
		zombieInterval,
		totalMessagePerPick,
	}

	return config
}

func (c *CommonOutboxConfig) GetGroupID() uuid.UUID {
	return c.GroupID
}

func (c *CommonOutboxConfig) GetOutboxTableName() string {
	return c.OutboxTableName
}

func (c *CommonOutboxConfig) GetPickPollInterval() int {
	return c.PickPollInterval
}

func (c *CommonOutboxConfig) GetPickZombiePollInterval() int {
	return c.PickZombiePollInterval
}

func (c *CommonOutboxConfig) GetZombieInterval() int {
	return c.ZombieInterval
}

func (c *CommonOutboxConfig) GetTotalMessagePerPick() int {
	return c.TotalMessagePerPick
}
