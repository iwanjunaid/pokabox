package manager

import (
	"database/sql"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/iwanjunaid/pokabox/event"
	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
	"github.com/iwanjunaid/pokabox/internal/interfaces/manager"
)

type CommonManager struct {
	outboxConfig  config.OutboxConfig
	kafkaConfig   config.KafkaConfig
	kafkaProducer *kafka.Producer
	eventHandler  event.EventHandler
	db            *sql.DB
	wg            *sync.WaitGroup
}

func New(outboxConfig config.OutboxConfig, kafkaConfig config.KafkaConfig, db *sql.DB) (manager.Manager, error) {
	var wg sync.WaitGroup

	kafkaProducer, err := kafka.NewProducer(kafkaConfig.GetConfigMap())

	if err != nil {
		return nil, err
	}

	manager := &CommonManager{
		outboxConfig:  outboxConfig,
		kafkaConfig:   kafkaConfig,
		kafkaProducer: kafkaProducer,
		eventHandler:  nil,
		db:            db,
		wg:            &wg,
	}

	return manager, nil
}

func (m *CommonManager) GetOutboxConfig() config.OutboxConfig {
	return m.outboxConfig
}

func (m *CommonManager) GetKafkaConfig() config.KafkaConfig {
	return m.kafkaConfig
}

func (m *CommonManager) GetKafkaProducer() *kafka.Producer {
	return m.kafkaProducer
}

func (m *CommonManager) SetEventHandler(e event.EventHandler) {
	m.eventHandler = e
}

func (m *CommonManager) GetEventHandler() event.EventHandler {
	return m.eventHandler
}

func (m *CommonManager) GetDB() *sql.DB {
	return m.db
}

func (m *CommonManager) Start() error {
	m.wg.Add(3)

	go backgroundPick(m)
	go backgroundZombiePick(m)
	go backgroundRemove(m)

	return nil
}

func (m *CommonManager) Await() {
	m.wg.Wait()
}
