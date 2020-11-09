package emitter

import (
	"time"

	"github.com/iwanjunaid/pokabox/model"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventSent(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID, kafkaTopic string, kafkaPartition int32,
	record *model.OutboxRecord) {
	if e != nil {
		eventSent := event.Sent{
			PickerGroupID:  pickerGroupID,
			KafkaTopic:     kafkaTopic,
			KafkaPartition: kafkaPartition,
			OutboxRecord:   record,
			Timestamp:      timestamp,
		}

		e(eventSent)
	}
}
