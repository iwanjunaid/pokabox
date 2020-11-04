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

func backgroundZombiePick(manager *CommonManager) {
	defer manager.wg.Done()

	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetZombiePickerPollInterval()
	messageLimit := outboxConfig.GetZombiePickerMessageLimitPerPoll()
	zombieInterval := outboxConfig.GetZombieInterval()
	tableName := outboxConfig.GetOutboxTableName()
	eventHandler := manager.GetEventHandler()

	q := `
	SELECT id, group_id, kafka_topic,
	kafka_key, kafka_value,
	priority, status, version,
	created_at, sent_at
	FROM %s
	WHERE status = 'NEW'
		AND group_id NOT IN ('%s')
		AND NOW() - created_at > '%d second'::interval
	ORDER BY priority, created_at
	LIMIT %d
	`

	query := fmt.Sprintf(q, tableName, outboxGroupID, zombieInterval, messageLimit)

	for {
		// Emit event ZombiePickerStarted
		if eventHandler != nil {
			eventZombiePickerStarted := event.ZombiePickerStarted{
				PickerGroupID: outboxGroupID,
				Timestamp:     time.Now(),
			}

			eventHandler(eventZombiePickerStarted)
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
				panic(err)
			}

			// Emit event ZombiePicked
			var record *model.OutboxRecord

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

				zombiePicked := event.ZombiePicked{
					PickerGroupID: outboxGroupID,
					OutboxRecord:  record,
					Timestamp:     time.Now(),
				}

				eventHandler(zombiePicked)
			}

			q := `
			UPDATE %s
			SET group_id = '%s',
				version = %d
			WHERE group_id = '%s'
				AND version = %d
			`

			newVersion := version.Int32 + 1
			updateQuery := fmt.Sprintf(q, tableName, outboxGroupID, newVersion, groupID.String, version.Int32)
			stmt, stmtErr := manager.GetDB().Prepare(updateQuery)

			if stmtErr != nil {
				log.Fatal(stmtErr)
			}

			_, execErr := stmt.Exec()

			if execErr != nil {
				log.Fatal(execErr)
			}

			// Emit event ZombieAcquired
			if eventHandler != nil {
				originGroupID := record.GroupID
				record.GroupID = outboxGroupID
				eventZombieAcquired := event.ZombieAcquired{
					PickerGroupID: outboxGroupID,
					OriginGroupID: originGroupID,
					OutboxRecord:  record,
					Timestamp:     time.Now(),
				}

				eventHandler(eventZombieAcquired)
			}

			rows.Close()
		}

		// Emit event ZombiePickerPaused
		if eventHandler != nil {
			eventZombiePickerPaused := event.ZombiePickerPaused{
				PickerGroupID: outboxGroupID,
				Timestamp:     time.Now(),
			}

			eventHandler(eventZombiePickerPaused)
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
