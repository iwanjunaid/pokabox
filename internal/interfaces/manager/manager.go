package manager

import (
	"database/sql"

	"github.com/iwanjunaid/pokabox/event"
	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
)

type Manager interface {
	GetOutboxConfig() config.OutboxConfig
	GetKafkaConfig() config.KafkaConfig
	GetDB() *sql.DB
	SetEventHandler(event.EventHandler)
	GetEventHandler() event.EventHandler
	// Insert(*sql.Tx, *model.OutboxRecord) error
	Start() error
	Await()
}
