package interfaces

import (
	"database/sql"

	"github.com/iwanjunaid/outbox-postgres/model"
)

type Manager interface {
	GetOutboxConfig() Config
	Insert(*sql.Tx, *model.OutboxRecord) error
	Start() error
	Await() error
}
