package manager

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/internal/interfaces/config"
	"github.com/iwanjunaid/pokabox/internal/interfaces/manager"
	"github.com/iwanjunaid/pokabox/model"
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

func backgroundPick(manager *CommonManager) {
	defer manager.wg.Done()

	fmt.Println("Picking...")
	outboxConfig := manager.GetOutboxConfig()
	groupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetPickPollInterval()
	tableName := outboxConfig.GetOutboxTableName()

	q := `
	SELECT id, group_id, kafka_topic,
		kafka_key, kafka_value,
		priority, status, version,
		created_at, sent_at
	FROM %s
	WHERE status = 'NEW'
		AND group_id = '%s'
	ORDER BY priority, created_at
	`

	query := fmt.Sprintf(q, tableName, groupID)

	for {
		rows, err := manager.GetDB().Query(query)

		if err != nil {
			panic(err)
		}

		var records []*model.OutboxRecord

		for rows.Next() {
			var (
				id         sql.NullString
				groupID    sql.NullString
				kafkaTopic sql.NullString
				kafkaKey   sql.NullString
				kafkaValue sql.NullString
				priority   sql.NullInt32
				status     sql.NullString
				version    sql.NullInt32
				createdAt  sql.NullTime
				sentAt     sql.NullTime
			)

			err := rows.Scan(&id, &groupID, &kafkaTopic, &kafkaKey,
				&kafkaValue, &priority, &status, &version,
				&createdAt, &sentAt)

			if err != nil {
				panic(err)
			}

			records = append(records, &model.OutboxRecord{
				ID:         uuid.MustParse(id.String),
				GroupID:    uuid.MustParse(groupID.String),
				KafkaTopic: kafkaTopic.String,
				KafkaKey:   kafkaKey.String,
				KafkaValue: kafkaValue.String,
				Priority:   uint(priority.Int32),
				Status:     status.String,
				Version:    uint(version.Int32),
				CreatedAt:  createdAt.Time,
				SentAt:     sentAt.Time,
			})

			// TODO: Send to kafka

			// Update status to 'SENT'
			q := `
			UPDATE %s
			SET status = 'SENT'
			WHERE id = '%s'
			`

			sentQuery := fmt.Sprintf(q, tableName, id.String)
			stmt, stmtErr := manager.GetDB().Prepare(sentQuery)

			if err != nil {
				log.Fatal(stmtErr)
			}

			_, execErr := stmt.Exec()

			if execErr != nil {
				log.Fatal(execErr)
			}

			rows.Close()
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func backgroundZombiePick(manager *CommonManager) {
	defer manager.wg.Done()

	fmt.Println("Picking zombies...")
	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetPickZombiePollInterval()
	zombieInterval := outboxConfig.GetZombieInterval()
	tableName := outboxConfig.GetOutboxTableName()

	query := fmt.Sprintf(`SELECT id, group_id, kafka_topic,
												kafka_key, kafka_value,
												priority, status, version,
												created_at, sent_at
												FROM %s
												WHERE status = 'NEW'
													AND group_id NOT IN ('%s')
													AND NOW() - created_at > '%d second'::interval
												ORDER BY priority, created_at`, tableName, outboxGroupID, zombieInterval)

	for {
		rows, err := manager.GetDB().Query(query)

		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var (
				id         sql.NullString
				groupID    sql.NullString
				kafkaTopic sql.NullString
				kafkaKey   sql.NullString
				kafkaValue sql.NullString
				priority   sql.NullInt32
				status     sql.NullString
				version    sql.NullInt32
				createdAt  sql.NullTime
				sentAt     sql.NullTime
			)

			err := rows.Scan(&id, &groupID, &kafkaTopic, &kafkaKey,
				&kafkaValue, &priority, &status, &version,
				&createdAt, &sentAt)

			if err != nil {
				panic(err)
			}

			newVersion := version.Int32 + 1
			updateQuery := fmt.Sprintf(`UPDATE %s
																SET group_id = '%s'
																	version = %d
																WHERE group_id = '%s'
																	AND version = %d`, tableName, outboxGroupID, newVersion, groupID.String, version.Int32)

			fmt.Println("***updateQuery", updateQuery)
			stmt, stmtErr := manager.GetDB().Prepare(updateQuery)

			if stmtErr != nil {
				log.Fatal(stmtErr)
			}

			_, execErr := stmt.Exec()

			if execErr != nil {
				log.Fatal(execErr)
			}

			rows.Close()
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}

}

func backgroundRemove(manager *CommonManager) {
	defer manager.wg.Done()

	fmt.Println("Removing sent records...")
	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	tableName := outboxConfig.GetOutboxTableName()

	query := fmt.Sprintf(`SELECT id
												FROM %s
												WHERE status = 'SENT'
													AND group_id = '%s'
												LIMIT 100`, tableName, outboxGroupID)

	for {
		rows, err := manager.GetDB().Query(query)

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var (
				id sql.NullString
			)

			err := rows.Scan(&id)

			if err != nil {
				log.Fatal(err)
			}

			deleteQuery := fmt.Sprintf(`DELETE FROM %s
																	WHERE status = 'SENT'
																		AND group_id = '%s'`, tableName, outboxGroupID)

			stmt, stmtErr := manager.GetDB().Prepare(deleteQuery)

			if stmtErr != nil {
				log.Fatal(stmtErr)
			}

			_, execErr := stmt.Exec()

			if execErr != nil {
				log.Fatal(execErr)
			}

			rows.Close()
		}
	}

	// TODO: sleep value should be fetch from config
	time.Sleep(4 * time.Second)
}
