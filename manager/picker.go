package manager

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event/emitter"
	"github.com/iwanjunaid/pokabox/model"
)

func backgroundPick(manager *CommonManager) {
	defer manager.wg.Done()
	defer func() {
		r := recover()

		fmt.Println("Recover from panic")

		if _, ok := r.(error); ok {
			manager.wg.Add(1)
			time.Sleep(3 * time.Second)

			go backgroundPick(manager)
		}
	}()

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
		// Emit event PickerStarted
		emitter.EmitEventPickerStarted(eventHandler, time.Now(), outboxGroupID)

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

			// Emit event Picked
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

			// Emit event Picked
			emitter.EmitEventPicked(eventHandler, time.Now(), record)

			// Send to kafka
			deliveryChan := make(chan kafka.Event, 10000)
			kafkaProducer := manager.GetKafkaProducer()

			var key []byte

			if kafkaKey.String != "" {
				key = []byte(kafkaKey.String)
			}

			errProduce := kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic.String, Partition: kafka.PartitionAny},
				Value:          []byte(kafkaValue.String),
				Key:            key,
			}, deliveryChan)

			if errProduce != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(), outboxGroupID, errProduce)
				panic(errProduce)
			}

			kafkaEvent := <-deliveryChan
			kafkaMessage := kafkaEvent.(*kafka.Message)

			if parErr := kafkaMessage.TopicPartition.Error; parErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(),
					outboxGroupID, parErr)
				panic(parErr)
			}

			// Emit event Sent
			emitter.EmitEventSent(eventHandler, time.Now(), outboxGroupID,
				*kafkaMessage.TopicPartition.Topic,
				kafkaMessage.TopicPartition.Partition,
				record)

			close(deliveryChan)

			// Update status to 'SENT'
			newSentAt := time.Now()
			q := `
			UPDATE %s
			SET status = $1,
				sent_at = $2
			WHERE id = '%s'
			`

			sentQuery := fmt.Sprintf(q, tableName, id.String)
			stmt, stmtErr := manager.GetDB().Prepare(sentQuery)

			if stmtErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(),
					outboxGroupID, stmtErr)
				panic(stmtErr)
			}

			_, execErr := stmt.Exec(model.FlagSent, newSentAt)

			if execErr != nil {
				emitter.EmitEventErrorOccured(eventHandler, time.Now(),
					outboxGroupID, execErr)
				panic(execErr)
			}

			// Emit event StatusChanged
			emitter.EmitEventStatusChanged(eventHandler, time.Now(),
				model.FlagNew, model.FlagSent, record)
		}

		rows.Close()

		if !recordsFound {
			// Emit event PickerPaused
			emitter.EmitEventPickerPaused(eventHandler, time.Now(), outboxGroupID)

			// Pause picker
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}
}
