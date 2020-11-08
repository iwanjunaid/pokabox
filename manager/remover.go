package manager

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event/emitter"
	"github.com/iwanjunaid/pokabox/model"
)

func backgroundRemove(manager *CommonManager) {
	defer manager.wg.Done()

	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetRemoverPollInterval()
	messageLimit := outboxConfig.GetRemoverMessageLimitPerPoll()
	tableName := outboxConfig.GetOutboxTableName()
	eventHandler := manager.GetEventHandler()

	q := `
	SELECT id, group_id, kafka_topic,
		kafka_key, kafka_value,
		priority, status, version,
		created_at, sent_at
	FROM %s
	WHERE status = 'SENT'
		AND group_id = '%s'
	LIMIT %d
	`

	query := fmt.Sprintf(q, tableName, outboxGroupID, messageLimit)

	for {
		// Emit event RemoverStarted
		emitter.EmitEventRemoverStarted(eventHandler, time.Now(), outboxGroupID)

		recordsFound := false
		rows, err := manager.GetDB().Query(query)

		if err != nil {
			emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, err)
			panic(err)
		}

		for rows.Next() {
			recordsFound = true

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
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, err)
				panic(err)
			}

			q := `
			DELETE FROM %s
			WHERE status = 'SENT' AND group_id = '%s'
			`

			deleteQuery := fmt.Sprintf(q, tableName, outboxGroupID)
			stmt, stmtErr := manager.GetDB().Prepare(deleteQuery)

			if stmtErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, stmtErr)
				panic(stmtErr)
			}

			_, execErr := stmt.Exec()

			if execErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, execErr)
				panic(execErr)
			}

			var record *model.OutboxRecord

			record = &model.OutboxRecord{
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
			}

			// Emit event Removed
			emitter.EmitEventRemoved(eventHandler, time.Now(), outboxGroupID, record)
		}

		rows.Close()

		if !recordsFound {
			// Emit event RemoverPaused
			emitter.EmitEventRemoverPaused(eventHandler, time.Now(), outboxGroupID)

			// Pause remover
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}
}
