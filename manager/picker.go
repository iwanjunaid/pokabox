package manager

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func backgroundPick(manager *CommonManager) {
	defer manager.wg.Done()

	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetPickerPollInterval()
	messageLimit := outboxConfig.GetPickerMessageLimitPerPoll()
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
	LIMIT %d
	`

	query := fmt.Sprintf(q, tableName, outboxGroupID, messageLimit)

	for {
		fmt.Println("Picking...")
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
