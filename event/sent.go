package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type Sent struct {
	PickerGroupID  uuid.UUID
	KafkaTopic     string
	KafkaPartition int32
	OutboxRecord   *model.OutboxRecord
	Timestamp      time.Time
}

func (s Sent) String() string {
	return fmt.Sprintf("[%s:%s] Message with ID %s sent to kafka in topic %s and partition %d",
		PREFIX, s.PickerGroupID, s.OutboxRecord.ID,
		s.KafkaTopic, s.KafkaPartition)
}

func (s Sent) GetPickerGroupID() uuid.UUID {
	return s.PickerGroupID
}

func (s Sent) GetRecord() *model.OutboxRecord {
	return s.OutboxRecord
}

func (s Sent) GetTimestamp() time.Time {
	return s.Timestamp
}
