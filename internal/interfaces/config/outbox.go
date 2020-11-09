package config

import "github.com/google/uuid"

type OutboxConfig interface {
	GetGroupID() uuid.UUID
	GetOutboxTableName() string
	GetPickerPollInterval() int
	GetPickerMessageLimitPerPoll() int
	GetZombieInterval() int
	GetZombiePickerPollInterval() int
	GetZombiePickerMessageLimitPerPoll() int
	GetRemoverPollInterval() int
	GetRemoverMessageLimitPerPoll() int
}
