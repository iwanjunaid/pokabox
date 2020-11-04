package manager

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
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
		if eventHandler != nil {
			eventRemoverStarted := event.RemoverStarted{
				PickerGroupID: outboxGroupID,
				Timestamp:     time.Now(),
			}

			eventHandler(eventRemoverStarted)
		}
		rows, err := manager.GetDB().Query(query)

		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}

			q := `
			DELETE FROM %s
			WHERE status = 'SENT'
				AND group_id = '%s'
			`

			deleteQuery := fmt.Sprintf(q, tableName, outboxGroupID)
			stmt, stmtErr := manager.GetDB().Prepare(deleteQuery)

			if stmtErr != nil {
				log.Fatal(stmtErr)
			}

			_, execErr := stmt.Exec()

			if execErr != nil {
				log.Fatal(execErr)
			}

			var record *model.OutboxRecord

			// Emit event Removed
			if eventHandler != nil {
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

				eventRemoved := event.Removed{
					PickerGroupID: outboxGroupID,
					OutboxRecord:  record,
					Timestamp:     time.Now(),
				}

				eventHandler(eventRemoved)
			}

			rows.Close()
		}

		// Emit event RemoverPaused
		if eventHandler != nil {
			eventRemoverPaused := event.RemoverPaused{
				PickerGroupID: outboxGroupID,
				Timestamp:     time.Now(),
			}

			eventHandler(eventRemoverPaused)
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
