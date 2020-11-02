package manager

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func backgroundZombiePick(manager *CommonManager) {
	defer manager.wg.Done()

	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	pollInterval := outboxConfig.GetZombiePickerPollInterval()
	messageLimit := outboxConfig.GetZombiePickerMessageLimitPerPoll()
	zombieInterval := outboxConfig.GetZombieInterval()
	tableName := outboxConfig.GetOutboxTableName()

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
		fmt.Println("Picking zombies...")
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

			q := `
			UPDATE %s
			SET group_id = '%s'
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

			rows.Close()
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
