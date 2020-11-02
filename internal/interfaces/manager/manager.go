package manager

import (
	"database/sql"

	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
)

type Manager interface {
	GetOutboxConfig() config.OutboxConfig
	GetKafkaConfig() config.KafkaConfig
	GetDB() *sql.DB
	// Insert(*sql.Tx, *model.OutboxRecord) error
	Start() error
	Await()
}
