package manager

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event/emitter"
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
		emitter.EmitEventZombiePickerStarted(eventHandler, time.Now(), outboxGroupID)

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

			// Emit event ZombiePicked
			emitter.EmitEventZombiePicked(eventHandler, time.Now(),
				outboxGroupID, record)

			q := `
			UPDATE %s
			SET group_id = '%s', version = %d
			WHERE id = '%s' AND version = %d
			`

			newVersion := version.Int32 + 1
			updateQuery := fmt.Sprintf(q, tableName, outboxGroupID, newVersion, id.String, version.Int32)
			stmt, stmtErr := manager.GetDB().Prepare(updateQuery)

			if stmtErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, stmtErr)
				panic(stmtErr)
			}

			result, execErr := stmt.Exec()

			if execErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, execErr)
				panic(execErr)
			}

			rowsAffected, affectedErr := result.RowsAffected()

			if affectedErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, affectedErr)
				panic(affectedErr)
			}

			// Emit event ZombieAcquired
			if rowsAffected > 0 {
				emitter.EmitEventZombieAcquired(eventHandler, time.Now(),
					outboxGroupID, record.GroupID, record)
			}
		}

		rows.Close()

		if !recordsFound {
			// Emit event ZombiePickerPaused
			emitter.EmitEventZombiePickerPaused(eventHandler, time.Now(), outboxGroupID)

			// Pause zombie picker
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}
}
