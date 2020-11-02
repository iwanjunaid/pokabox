package model

import (
	"time"

	"github.com/google/uuid"
)

type OutboxRecord struct {
	ID         uuid.UUID
	GroupID    uuid.UUID
	KafkaTopic string
	KafkaKey   string
	KafkaValue string
	Priority   uint
	Status     string
	Version    uint
	CreatedAt  time.Time
	SentAt     time.Time
}
