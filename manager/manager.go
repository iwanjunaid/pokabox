package manager

import (
	"database/sql"
	"sync"

	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
	"github.com/iwanjunaid/pokabox/internal/interfaces/manager"
)

type CommonManager struct {
	outboxConfig config.OutboxConfig
	kafkaConfig  config.KafkaConfig
	db           *sql.DB
	wg           sync.WaitGroup
}

func New(outboxConfig config.OutboxConfig, kafkaConfig config.KafkaConfig, db *sql.DB) manager.Manager {
	var wg sync.WaitGroup

	manager := &CommonManager{
		outboxConfig,
		kafkaConfig,
		db,
		wg,
	}

	return manager
}

func (m *CommonManager) GetOutboxConfig() config.OutboxConfig {
	return m.outboxConfig
}

func (m *CommonManager) GetKafkaConfig() config.KafkaConfig {
	return m.kafkaConfig
}

func (m *CommonManager) GetDB() *sql.DB {
	return m.db
}

func (m *CommonManager) Start() error {
	m.wg.Add(1)
	go backgroundPick(m)

	m.wg.Add(1)
	go backgroundZombiePick(m)

	m.wg.Add(1)
	go backgroundRemove(m)

	return nil
}

func (m *CommonManager) Await() {
	m.wg.Wait()
}
