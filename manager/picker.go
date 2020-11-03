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

func backgroundPick(manager *CommonManager) {
	defer manager.wg.Done()

	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetPickerPollInterval()
	messageLimit := outboxConfig.GetPickerMessageLimitPerPoll()
	tableName := outboxConfig.GetOutboxTableName()
	eventHandler := manager.GetEventHandler()

	q := `
	SELECT id, group_id, kafka_topic,
		kafka_key, kafka_value,
		priority, status, version,
		created_at, sent_at
	FROM %s
	WHERE status = '%s'
		AND group_id = '%s'
	ORDER BY priority, created_at
	LIMIT %d
	`

	query := fmt.Sprintf(q, tableName, model.FlagNew, outboxGroupID, messageLimit)

	for {
		if eventHandler != nil {
			pickerStarted := event.PickerStarted{
				GroupID: outboxGroupID,
			}

			eventHandler(pickerStarted)
		}

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

			if eventHandler != nil {
				record := &model.OutboxRecord{
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

				fetched := event.Fetched{
					OutboxRecord: record,
				}

				eventHandler(fetched)
			}

			// TODO: Send to kafka

			// Update status to 'SENT'
			q := `
			UPDATE %s
			SET status = '%s'
			WHERE id = '%s'
			`

			sentQuery := fmt.Sprintf(q, tableName, model.FlagSent, id.String)
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

		if eventHandler != nil {
			pickerPaused := event.PickerPaused{
				GroupID: outboxGroupID,
			}

			eventHandler(pickerPaused)
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
