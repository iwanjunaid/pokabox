package manager

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func backgroundRemove(manager *CommonManager) {
	defer manager.wg.Done()

	outboxConfig := manager.GetOutboxConfig()
	outboxGroupID := outboxConfig.GetGroupID()
	tableName := outboxConfig.GetOutboxTableName()

	q := `
	SELECT id
	FROM %s
	WHERE status = 'SENT'
		AND group_id = '%s'
	LIMIT 100
	`

	query := fmt.Sprintf(q, tableName, outboxGroupID)

	for {
		fmt.Println("Removing sent records...")
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

			rows.Close()
		}

		// TODO: sleep value should be fetch from config
		time.Sleep(4 * time.Second)
	}
}
