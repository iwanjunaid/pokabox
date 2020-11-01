package model

import (
	"github.com/google/uuid"
	"time"
)

type OutboxRecord struct {
	ID         uuid.UUID
	GroupID    uuid.UUID
	KafkaTopic string
	KafkaKey   string
	KafkaValue string
	Priority   uint
	Status     string
	CreatedAt  time.Time
	SentAt     time.Time
}
