package config

import "github.com/google/uuid"

type OutboxConfig interface {
	GetGroupID() uuid.UUID
	GetOutboxTableName() string
	GetPickPollInterval() int
	GetPickZombiePollInterval() int
	GetZombieInterval() int
	GetTotalMessagePerPick() int
}
